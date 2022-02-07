/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runtime

import (
	"context"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

var logger = ctrl.Log.WithName("runtime").WithName("reconcile")

type (
	Reconciler interface {
		Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error)
		Next(ctx context.Context, obj client.Object) (ctrl.Result, error)
		setNext(next Reconciler)
	}
	ReconcilerFuncs struct {
		ReconcilerResult
		next Reconciler
	}
	ReconcilerResult struct{}
)

func Reconcilers(reconcilers ...Reconciler) Reconciler {
	reconcilers = append(reconcilers, &finisher{})
	reconcilers = append([]Reconciler{&trace{}}, reconcilers...)
	var last Reconciler
	for i := len(reconcilers) - 1; i >= 0; i-- {
		current := reconcilers[i]
		current.setNext(last)
		last = current
	}
	return last
}

func (r *ReconcilerFuncs) setNext(next Reconciler) {
	r.next = next
}

func (r *ReconcilerFuncs) Next(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	return r.next.Reconcile(ctx, obj)
}

func (r *ReconcilerResult) Finish() (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (r *ReconcilerResult) Requeue() (ctrl.Result, error) {
	return ctrl.Result{Requeue: true}, nil
}

func (r *ReconcilerResult) RequeueAfter(duration time.Duration) (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: duration}, nil
}

func (r *ReconcilerResult) RequeueOnErr(ctx context.Context, err error) (ctrl.Result, error) {
	if span := tracing.SpanFromContext(ctx); span != nil && err != nil {
		span.SetError(err)
	}
	return ctrl.Result{}, err
}

func (r *ReconcilerResult) RequeueOnErrAfter(ctx context.Context, err error, duration time.Duration) (ctrl.Result, error) {
	if span := tracing.SpanFromContext(ctx); span != nil && err != nil {
		span.SetError(err)
	}
	return ctrl.Result{RequeueAfter: duration}, err
}
