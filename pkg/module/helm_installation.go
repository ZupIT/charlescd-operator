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
	"github.com/fluxcd/pkg/apis/meta"
	"github.com/fluxcd/source-controller/api/v1beta1"
	mf "github.com/manifestival/manifestival"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

type Helm interface {
	Template(name, chart string, values map[string]interface{}) (mf.Manifest, error)
}

type GitRepositoryGetter interface {
	GetGitRepository(ctx context.Context, key client.ObjectKey) (*v1beta1.GitRepository, error)
}

type HelmInstallation struct {
	reconciler.Funcs

	git    GitRepositoryGetter
	status StatusWriter
}

func NewHelmInstallation(git GitRepositoryGetter, status StatusWriter) *HelmInstallation {
	return &HelmInstallation{git: git, status: status}
}

func (h *HelmInstallation) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	if module, ok := obj.(*deployv1alpha1.Module); ok {
		return h.reconcile(ctx, module)
	}
	return h.Next(ctx, obj)
}

func (h *HelmInstallation) reconcile(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	l := logger.WithValues("trace", span)

	repo, err := h.git.GetGitRepository(ctx, types.NamespacedName{
		Namespace: module.GetNamespace(),
		Name:      module.GetName(),
	})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return h.Next(ctx, module)
		}
		l.Error(err, "Error getting git repository")
		return h.RequeueOnErr(ctx, err)
	}

	if c := apimeta.FindStatusCondition(repo.Status.Conditions, meta.ReadyCondition); c != nil {
		if c.Status == metav1.ConditionFalse {
			if diff, updated := module.SetSourceError("GitRepositoryError", c.Message); updated {
				l.Info("Status changed", "diff", diff)
				return h.RequeueOnErr(ctx, h.status.UpdateModuleStatus(ctx, module))
			}
		}
	} else {
		if diff, updated := module.RemoveSource(); updated {
			l.Info("Status changed", "diff", diff)
			return h.RequeueOnErr(ctx, h.status.UpdateModuleStatus(ctx, module))
		}
	}

	artifact := repo.GetArtifact()
	if artifact == nil {
		l.Info("The artifact is not ready")
		return h.Next(ctx, module)
	}

	u, err := url.Parse(artifact.URL)
	if err != nil {
		l.Error(err, "Error reading artifact address")
		if diff, updated := module.SetSourceError("AddressResolutionError", err.Error()); updated {
			l.Info("Status changed", "diff", diff)
			return h.RequeueOnErr(ctx, h.status.UpdateModuleStatus(ctx, module))
		}
		return h.RequeueOnErr(ctx, err)
	}

	filepath := os.TempDir() + u.Path
	if err = downloadArtifact(ctx, filepath, artifact); err != nil {
		l.Error(err, "Error downloading artifact")
		if diff, updated := module.SetSourceError("DownloadError", err.Error()); updated {
			l.Info("Status changed", "diff", diff)
			return h.RequeueOnErr(ctx, h.status.UpdateModuleStatus(ctx, module))
		}
		return h.RequeueOnErr(ctx, err)
	}

	if diff, updated := module.SetSourceReady(filepath); updated {
		l.Info("Status changed", "diff", diff)
		return h.RequeueOnErr(ctx, h.status.UpdateModuleStatus(ctx, module))
	}

	//	TODO: implement Helm client
	//	manifest, err := h.Helm.Template(module.GetName(), filepath, module.Spec.Values)
	//	if err != nil {
	//		l.Error(err, "Error rendering Helm chart templates")
	//		return runtime.Finish()
	//	}
	//
	//	if err = manifest.Apply(); err != nil {
	//		l.Error(err, "Error applying Helm chart changes")
	//		return runtime.RequeueOnErr(ctx, err)
	//	}

	return h.Next(ctx, module)
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
	// url := "http://127.0.0.1:9090/gitrepository/default/football-bets/da684f367e901b0e2747a69c2914bd9382b1428e.tar.gz"
	// res, err := http.Get(url)
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