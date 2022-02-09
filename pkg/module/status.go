package module

import (
	"context"

	"github.com/angelokurtis/reconciler"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

type Status struct {
	reconciler.Funcs
	status StatusWriter
}

func NewStatus(status StatusWriter) *Status {
	return &Status{status: status}
}

func (s *Status) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	if module, ok := obj.(*deployv1alpha1.Module); ok {
		return s.reconcile(ctx, module)
	}
	return s.Next(ctx, obj)
}

func (s *Status) reconcile(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	l := logger.WithValues("trace", span)

	if diff, updated := module.UpdatePhase(); updated {
		l.Info("Status changed", "diff", diff)
		return s.RequeueOnErr(ctx, s.status.UpdateModuleStatus(ctx, module))
	}
	return s.Next(ctx, module)
}
