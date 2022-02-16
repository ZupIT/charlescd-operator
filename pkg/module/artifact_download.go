/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package module

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/angelokurtis/reconciler"
	"github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/go-logr/logr"
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

const (
	addressResolutionError = "AddressResolutionError"
	downloadError          = "DownloadError"
	gitRepositoryError     = "GitRepositoryError"
)

type Helm interface {
	Template(name, chart string, values map[string]interface{}) (mf.Manifest, error)
}

type GitRepositoryGetter interface {
	GetGitRepository(ctx context.Context, key client.ObjectKey) (*v1beta1.GitRepository, error)
}

type ArtifactDownload struct {
	reconciler.Funcs

	git    GitRepositoryGetter
	status StatusWriter
}

func NewArtifactDownload(git GitRepositoryGetter, status StatusWriter) *ArtifactDownload {
	return &ArtifactDownload{git: git, status: status}
}

func (a *ArtifactDownload) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	if module, ok := obj.(*deployv1alpha1.Module); ok && module.Spec.Repository.Git != nil {
		return a.reconcile(ctx, module)
	}
	return a.Next(ctx, obj)
}

func (a *ArtifactDownload) reconcile(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := logr.FromContextOrDiscard(ctx)

	repo, err := a.git.GetGitRepository(ctx, types.NamespacedName{
		Namespace: module.GetNamespace(),
		Name:      module.GetName(),
	})
	if err != nil {
		return a.RequeueOnErr(ctx, err)
	}

	// check if this handler should act
	if repo == nil && repo.GetArtifact() == nil {
		log.Info("Artifact is not ready")
		return a.Next(ctx, module)
	}

	// get artifact address
	artifact := repo.GetArtifact()
	if msg, ok := statusOf(repo).IsError(); ok && module.SetSourceError(gitRepositoryError, msg) {
		return a.status.UpdateModuleStatus(ctx, module)
	}
	u, err := url.Parse(artifact.URL)
	if err != nil {
		log.Error(err, "Error reading artifact address")
		if module.SetSourceError(addressResolutionError, err.Error()) {
			return a.status.UpdateModuleStatus(ctx, module)
		}
		return a.RequeueOnErr(ctx, err)
	}
	filepath := os.TempDir() + u.Path

	// search for artifact locally
	if _, err = os.Stat(filepath); !errors.Is(err, os.ErrNotExist) && a.checksum(filepath, artifact.Checksum) {
		log.Info("Artifact found locally", "path", filepath, "checksum", artifact.Checksum)
		return a.updateStatusToReady(ctx, module, filepath)
	}

	// download the artifact
	if err = a.download(ctx, filepath, artifact); err != nil {
		log.Error(err, "Error downloading artifact")
		if module.SetSourceError(downloadError, err.Error()) {
			return a.status.UpdateModuleStatus(ctx, module)
		}
		return a.RequeueOnErr(ctx, err)
	}

	log.Info("Artifact downloaded", "path", filepath, "checksum", artifact.Checksum)
	return a.updateStatusToReady(ctx, module, filepath)
}

func (a *ArtifactDownload) download(ctx context.Context, filepath string, artifact *v1beta1.Artifact) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	// u := "http://127.0.0.1:9090/gitrepository/default/football-bets/da684f367e901b0e2747a69c2914bd9382b1428e.tar.gz"
	// res, err := http.Get(u)
	res, err := http.Get(artifact.URL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	index := strings.LastIndex(filepath, "/")
	if err = os.MkdirAll(filepath[:index], os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, res.Body); err != nil {
		return err
	}

	return nil
}

func (a *ArtifactDownload) checksum(filepath, checksum string) bool {
	f, err := os.Open(filepath)
	if err != nil {
		return false
	}
	defer f.Close()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return false
	}

	return fmt.Sprintf("%x", h.Sum(nil)) == checksum
}

func (a *ArtifactDownload) updateStatusToReady(ctx context.Context, module *deployv1alpha1.Module, filepath string) (ctrl.Result, error) {
	if module.SetSourceReady(filepath) {
		return a.status.UpdateModuleStatus(ctx, module)
	}

	return a.Next(ctx, module)
}
