package runtime

import (
	"context"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

func Finish() (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func Requeue() (ctrl.Result, error) {
	return ctrl.Result{Requeue: true}, nil
}

func RequeueAfter(duration time.Duration) (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: duration}, nil
}

func RequeueOnErr(ctx context.Context, err error) (ctrl.Result, error) {
	if span := tracing.SpanFromContext(ctx); span != nil {
		span.SetError(err)
	}
	return ctrl.Result{}, err
}

func RequeueOnErrAfter(ctx context.Context, err error, duration time.Duration) (ctrl.Result, error) {
	if span := tracing.SpanFromContext(ctx); span != nil {
		span.SetError(err)
	}
	return ctrl.Result{RequeueAfter: duration}, err
}
