package module

import (
	"context"

	"github.com/angelokurtis/reconciler"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

type HelmValidation struct {
	reconciler.Funcs
}

func (h *HelmValidation) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	//	TODO: implement Helm client
	//	manifest, err := h.Helm.Template(module.GetName(), filepath, module.Spec.Values)
	//	if err != nil {
	//		log.Error(err, "Error rendering Helm chart templates")
	//		return runtime.Finish()
	//	}
	//
	//	if err = manifest.Apply(); err != nil {
	//		log.Error(err, "Error applying Helm chart changes")
	//		return runtime.RequeueOnErr(ctx, err)
	//	}

	return h.Next(ctx, obj)
}
