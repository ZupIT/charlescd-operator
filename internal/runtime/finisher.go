package runtime

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type finisher struct{ ReconcilerFuncs }

func (f *finisher) Reconcile(context.Context, client.Object) (ctrl.Result, error) {
	return f.Finish()
}
