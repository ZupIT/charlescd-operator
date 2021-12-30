package usecase

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

	"github.com/fluxcd/source-controller/api/v1beta1"
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/runtime"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

type Helm interface {
	Template(name, chart string, values map[string]interface{}) (mf.Manifest, error)
}

type GitRepositoryGetter interface {
	GetGitRepository(ctx context.Context, key client.ObjectKey) (*v1beta1.GitRepository, error)
}

type HelmInstallation struct {
	GitRepositoryGetter GitRepositoryGetter
	Helm                Helm
}

func (hi *HelmInstallation) EnsureHelmInstallation(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	l := logger.WithValues("trace", span)

	repo, err := hi.GitRepositoryGetter.GetGitRepository(ctx, types.NamespacedName{
		Namespace: module.GetNamespace(),
		Name:      module.GetName(),
	})
	if err != nil {
		l.Error(err, "Error getting git repository")
		return runtime.RequeueOnErr(ctx, err)
	}

	artifact := repo.GetArtifact()
	if artifact == nil {
		l.Info("The artifact is not ready")
		return runtime.Finish()
	}

	u, err := url.Parse(artifact.URL)
	if err != nil {
		l.Error(err, "Error reading artifact address")
		return runtime.RequeueOnErr(ctx, err)
	}

	filepath := "." + u.Path
	if err = downloadArtifact(ctx, filepath, artifact); err != nil {
		l.Error(err, "Error downloading artifact")
		return runtime.RequeueOnErr(ctx, err)
	}

	//	TODO: implement Helm client
	//	manifest, err := hi.Helm.Template(module.GetName(), filepath, module.Spec.Values)
	//	if err != nil {
	//		l.Error(err, "Error rendering Helm chart templates")
	//		return runtime.Finish()
	//	}
	//
	//	if err = manifest.Apply(); err != nil {
	//		l.Error(err, "Error applying Helm chart changes")
	//		return runtime.RequeueOnErr(ctx, err)
	//	}

	return runtime.Finish()
}

func downloadArtifact(ctx context.Context, filepath string, artifact *v1beta1.Artifact) error {
	span := tracing.SpanFromContext(ctx)
	l := logger.WithValues("trace", span)

	if _, err := os.Stat(filepath); !errors.Is(err, os.ErrNotExist) && checksumIsValid(filepath, artifact.Checksum) {
		l.Info("Artifact found locally", "path", filepath, "checksum", artifact.Checksum)
		return nil
	}

	span, ctx = tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	l = logger.WithValues("trace", span)
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
	l.Info("Downloaded artifact", "path", filepath, "checksum", artifact.Checksum)
	return nil
}

func checksumIsValid(filepath, checksum string) bool {
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
