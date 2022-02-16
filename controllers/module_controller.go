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

package controllers

import (
	"context"

	"github.com/angelokurtis/reconciler"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/event"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
	"github.com/tiagoangelozup/charles-alpha/pkg/module"
)

type ModuleGetter interface {
	GetModule(ctx context.Context, key client.ObjectKey) (*deployv1alpha1.Module, error)
}

type ModuleHandler struct {
	*module.Status
	*module.DesiredState
	*module.ArtifactDownload
	*module.HelmValidation
}

// ModuleReconciler reconciles a Module object
type ModuleReconciler struct {
	reconciler.Result
	handle *ModuleHandler
	client ModuleGetter
}

func newModuleReconciler(handle *ModuleHandler, client ModuleGetter) *ModuleReconciler {
	return &ModuleReconciler{handle: handle, client: client}
}

//+kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=gitrepositories,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=deploy.charlescd.io,resources=modules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=deploy.charlescd.io,resources=modules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=deploy.charlescd.io,resources=modules/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ModuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := logr.FromContextOrDiscard(ctx)
	ctx = logf.IntoContext(ctx, log.WithValues("name", req.Name, "namespace", req.Namespace))

	m, err := r.client.GetModule(ctx, req.NamespacedName)
	if err != nil {
		log.Error(err, "Error getting resource with desired module state")
		return r.RequeueOnErr(ctx, err)
	}
	if m == nil {
		return r.Finish(ctx) // Module resource not found. Ignoring since object must be deleted
	}

	ctx = logf.IntoContext(ctx, log.WithValues("name", m.Name, "namespace", m.Namespace, "resourceVersion", m.ResourceVersion))
	return reconciler.Chain(
		r.handle.Status,
		r.handle.DesiredState,
		r.handle.ArtifactDownload,
		r.handle.HelmValidation,
	).Reconcile(ctx, m)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ModuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&deployv1alpha1.Module{}).
		Owns(&sourcev1.GitRepository{}).
		WithEventFilter(predicate.Or(
			event.NewRepoStatusPredicate(),
			event.NewModulePredicate(),
		)).
		WithLogger(logr.Discard()).
		Complete(r)
}
