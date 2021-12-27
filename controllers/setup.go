/*
Copyright 2021.

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

package controllers

import (
	"context"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func SetupWithManager(mgr ctrl.Manager) error {
	client, scheme := mgr.GetClient(), mgr.GetScheme()
	r, err := createReconcilers(client, scheme)
	if err != nil {
		return fmt.Errorf("unable to create reconcilers: %w", err)
	}
	for _, reconciler := range r {
		if err = reconciler.SetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to setup %T with manager: %w", reconciler, err)
		}
	}
	return nil
}

type Reconciler interface {
	Reconcile(context.Context, reconcile.Request) (reconcile.Result, error)
	SetupWithManager(ctrl.Manager) error
}
