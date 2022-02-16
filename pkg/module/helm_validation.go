package module

import (
	"context"

	"github.com/angelokurtis/reconciler"
	"github.com/go-logr/logr"
	mf "github.com/manifestival/manifestival"
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
		Template(ctx context.Context, releaseName, source, path string) (mf.Manifest, error)
	}
	HelmValidation struct {
		reconciler.Funcs

		helm   HelmClient
		status StatusWriter
	}
)

func NewHelmValidation(helm HelmClient, status StatusWriter) *HelmValidation {
	return &HelmValidation{helm: helm, status: status}
}

func (h *HelmValidation) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	if module, ok := obj.(*deployv1alpha1.Module); ok {
		return h.reconcile(ctx, module)
	}
	return h.Next(ctx, obj)
}

func (h *HelmValidation) reconcile(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error) {
	// check if this handler should act
	if module.Status.Source == nil || module.Status.Source.Path == "" {
		return h.Next(ctx, module)
	}

	// starting the context
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := logr.FromContextOrDiscard(ctx)

	source, path := module.Status.Source.Path, ""
	if module.Spec.Repository.Git != nil && module.Spec.Repository.Git.Path != "" {
		path = module.Spec.Repository.Git.Path
	}

	// templating Helm chart
	manifests, err := h.helm.Template(ctx, module.GetName(), source, path)
	if err != nil {
		log.Error(err, "Error templating source")
		if diff, updated := module.SetSourceInvalid(renderError, err.Error()); updated {
			log.Info("Status changed", "diff", diff)
			return h.RequeueOnErr(ctx, h.status.UpdateModuleStatus(ctx, module))
		}
		return h.Next(ctx, module)
	}

	// validate Helm chart templates
	if _, err = manifests.DryRun(); err != nil {
		log.Error(err, "Error validating Helm chart templates")
		if diff, updated := module.SetSourceInvalid(kubernetesAPIError, err.Error()); updated {
			log.Info("Status changed", "diff", diff)
			return h.RequeueOnErr(ctx, h.status.UpdateModuleStatus(ctx, module))
		}
		return h.Next(ctx, module)
	}

	// update status to success
	if diff, updated := module.SetSourceValid(); updated {
		log.Info("Status changed", "diff", diff)
		return h.RequeueOnErr(ctx, h.status.UpdateModuleStatus(ctx, module))
	}

	log.Info("Helm chart is valid")
	return h.Next(ctx, module)
}
