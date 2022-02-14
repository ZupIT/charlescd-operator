package module

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/angelokurtis/reconciler"
	"github.com/hashicorp/go-getter"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

type HelmValidation struct {
	reconciler.Funcs
}

func NewHelmValidation() *HelmValidation { return &HelmValidation{} }

func (h *HelmValidation) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	if module, ok := obj.(*deployv1alpha1.Module); ok {
		return h.reconcile(ctx, module)
	}
	return h.Next(ctx, obj)
}

func (h *HelmValidation) reconcile(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var origin, destination string
	if module.Status.Source != nil && module.Status.Source.Path != "" {
		origin = module.Status.Source.Path
		destination = origin[0 : len(origin)-len(".tar.gz")]
	}

	if module.Spec.Repository.Git != nil && module.Spec.Repository.Git.Path != "" {
		origin += "//" + module.Spec.Repository.Git.Path
		destination = filepath.Join(destination, module.Spec.Repository.Git.Path)
	}

	if err := getter.GetAny(destination, origin); err != nil {
		return h.RequeueOnErr(ctx, fmt.Errorf("error extracting Source artifact: %w", err))
	}

	//	TODO: implement Helm client
	//manifest, err := h.Helm.Template(module.GetName(), filepath, module.Spec.Values)
	//if err != nil {
	//	log.Error(err, "Error rendering Helm chart templates")
	//	return runtime.Finish()
	//}
	//
	//if err = manifest.Apply(); err != nil {
	//	log.Error(err, "Error applying Helm chart changes")
	//	return runtime.RequeueOnErr(ctx, err)
	//}

	return h.Next(ctx, module)
}
