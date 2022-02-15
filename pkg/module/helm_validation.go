package module

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/angelokurtis/reconciler"
	"github.com/go-logr/logr"
	"github.com/hashicorp/go-getter"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

const (
	kubernetesAPIError = "KubernetesAPIError"
	renderError        = "RenderError"
)

type (
	HelmClient interface {
		Template(ctx context.Context, name, chart string, values runtime.RawExtension) (*release.Release, error)
	}
	HelmValidation struct {
		reconciler.Funcs

		helm     HelmClient
		manifest ManifestReader
		status   StatusWriter
	}
)

func NewHelmValidation(helm HelmClient, manifest ManifestReader, status StatusWriter) *HelmValidation {
	return &HelmValidation{helm: helm, manifest: manifest, status: status}
}

func (h *HelmValidation) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	if module, ok := obj.(*deployv1alpha1.Module); ok {
		return h.reconcile(ctx, module)
	}
	return h.Next(ctx, obj)
}

func (h *HelmValidation) reconcile(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := logr.FromContextOrDiscard(ctx)

	if module.Status.Source == nil || module.Status.Source.Path == "" {
		return h.Next(ctx, module)
	}

	origin := module.Status.Source.Path
	destination := origin[0 : len(origin)-len(".tar.gz")]

	if module.Spec.Repository.Git != nil && module.Spec.Repository.Git.Path != "" {
		origin += "//" + module.Spec.Repository.Git.Path
		destination = filepath.Join(destination, module.Spec.Repository.Git.Path)
	}

	if err := getter.GetAny(destination, origin); err != nil {
		return h.RequeueOnErr(ctx, fmt.Errorf("error extracting Source artifact: %w", err))
	}

	rel, err := h.helm.Template(ctx, module.GetName(), destination, module.Spec.Values)
	if err != nil {
		log.Error(err, "Error rendering Helm chart templates")
		if diff, updated := module.SetSourceInvalid(renderError, err.Error()); updated {
			log.Info("Status changed", "diff", diff)
			return h.RequeueOnErr(ctx, h.status.UpdateModuleStatus(ctx, module))
		}
		return h.Next(ctx, module)
	}

	manifest, err := h.manifest.FromString(ctx, rel.Manifest)
	if err != nil {
		log.Error(err, "Error reading rendered Helm chart templates")
		if diff, updated := module.SetSourceInvalid(renderError, err.Error()); updated {
			log.Info("Status changed", "diff", diff)
			return h.RequeueOnErr(ctx, h.status.UpdateModuleStatus(ctx, module))
		}
		return h.Next(ctx, module)
	}

	if _, err = manifest.DryRun(); err != nil {
		log.Error(err, "Error validating Helm chart templates")
		if diff, updated := module.SetSourceInvalid(kubernetesAPIError, err.Error()); updated {
			log.Info("Status changed", "diff", diff)
			return h.RequeueOnErr(ctx, h.status.UpdateModuleStatus(ctx, module))
		}
		return h.Next(ctx, module)
	}

	if diff, updated := module.SetSourceValid(rel.Manifest); updated {
		log.Info("Status changed", "diff", diff)
		return h.RequeueOnErr(ctx, h.status.UpdateModuleStatus(ctx, module))
	}

	log.Info("Helm chart is valid")
	return h.Next(ctx, module)
}
