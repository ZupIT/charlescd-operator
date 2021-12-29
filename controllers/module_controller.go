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

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/runtime"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

var logger = ctrl.Log.WithName("controller").WithName("module")

type ModuleGetter interface {
	GetModule(ctx context.Context, key client.ObjectKey) (*deployv1alpha1.Module, error)
}

// ModuleReconciler reconciles a Module object
type ModuleReconciler struct {
	ModuleAdapter
	ModuleGetter ModuleGetter
}

//+kubebuilder:rbac:groups=deploy.charlescd.io,resources=modules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=deploy.charlescd.io,resources=modules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=deploy.charlescd.io,resources=modules/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ModuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	l := logger.WithValues("trace", span)

	l.Info("Reconciling...")
	module, err := r.ModuleGetter.GetModule(ctx, req.NamespacedName)
	if err != nil {
		l.Error(err, "Error getting resource with desired module state")
		return runtime.RequeueOnErr(ctx, err)
	}

	return NewModuleOperationChain(
		r.EnsureDesiredState,
	).Handle(ctx, module)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ModuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&deployv1alpha1.Module{}).
		Watches(
			&source.Kind{Type: &sourcev1.GitRepository{}},
			&handler.EnqueueRequestForOwner{OwnerType: &deployv1alpha1.Module{}, IsController: false}).
		WithLogger(log.NullLogger{}).
		Complete(r)
}
