package runtime

import (
	"context"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

type trace struct{ ReconcilerFuncs }

func (t *trace) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	span := tracing.SpanFromContext(ctx)
	l := logger.WithValues("trace", span)
	l.Info("Reconciler has been triggered")
	result, err := t.next.Reconcile(ctx, obj)
	switch {
	case err != nil:
		l.Error(err, "Reconciler error")
		return result, err
	case result.RequeueAfter > 0:
		l.Info("Successfully reconciled!", "requeue", fmt.Sprintf("in %s", result.RequeueAfter))
		return result, nil
	case result.Requeue:
		l.Info("Successfully reconciled!", "requeue", "now")
		return result, nil
	}
	l.Info("Successfully reconciled!")
	return result, nil
}
