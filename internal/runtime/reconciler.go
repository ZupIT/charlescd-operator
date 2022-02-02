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
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

var logger = ctrl.Log.WithName("runtime").WithName("reconcile")

type ReconcilerOperation interface {
	Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error)
	SetNext(next ReconcilerOperation)
}

func Operations(operations ...ReconcilerOperation) ReconcilerOperation {
	operations = append(operations, &doNothing{})
	operations = append([]ReconcilerOperation{&trace{}}, operations...)
	var last ReconcilerOperation
	for i := len(operations) - 1; i >= 0; i-- {
		current := operations[i]
		current.SetNext(last)
		last = current
	}
	return last
}

type trace struct{ next ReconcilerOperation }

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

func (t *trace) SetNext(next ReconcilerOperation) {
	t.next = next
}

type doNothing struct{}

func (d *doNothing) Reconcile(context.Context, client.Object) (ctrl.Result, error) {
	return Finish()
}

func (d *doNothing) SetNext(ReconcilerOperation) {}
