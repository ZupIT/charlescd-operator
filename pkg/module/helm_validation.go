package module

import (
	"context"

	"github.com/angelokurtis/reconciler"
	"github.com/go-logr/logr"
	mf "github.com/manifestival/manifestival"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
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
		Template(ctx context.Context, releaseName, source, path string, values *apiextensionsv1.JSON) (mf.Manifest, error)
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
	manifests, err := h.helm.Template(ctx, module.GetName(), source, path, module.Spec.Values)
	if err != nil {
		log.Error(err, "Error templating source")
		if module.SetSourceInvalid(renderError, err.Error()) {
			return h.status.UpdateModuleStatus(ctx, module)
		}
		return h.Next(ctx, module)
	}

	// validate Helm chart templates
	if _, err = manifests.DryRun(); err != nil {
		log.Error(err, "Error validating Helm chart templates")
		if module.SetSourceInvalid(kubernetesAPIError, err.Error()) {
			return h.status.UpdateModuleStatus(ctx, module)
		}
		return h.Next(ctx, module)
	}

	// update status to success
	if module.SetSourceValid() {
		return h.status.UpdateModuleStatus(ctx, module)
	}

	log.Info("Helm chart is valid")
	return h.Next(contextWithResources(ctx, manifests.Resources()), module)
}
