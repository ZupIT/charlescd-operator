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
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
)

var logger = ctrl.Log.WithName("controller").WithName("module")

type Manifests interface {
	Defaults(ctx context.Context) (mf.Manifest, error)
}

type ModuleGetter interface {
	GetModule(ctx context.Context, key client.ObjectKey) (*deployv1alpha1.Module, error)
}

// ModuleReconciler reconciles a Module object
type ModuleReconciler struct {
	Manifests    Manifests
	ModuleGetter ModuleGetter
	Scheme       *runtime.Scheme
}

//+kubebuilder:rbac:groups=deploy.charlescd.io,resources=modules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=deploy.charlescd.io,resources=modules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=deploy.charlescd.io,resources=modules/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ModuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := logger.V(1)

	l.Info("Reconciling...")
	module, err := r.ModuleGetter.GetModule(ctx, req.NamespacedName)
	if err != nil {
		l.Error(err, "Error getting resource with desired module state")
		return ctrl.Result{}, err
	}

	var manifests mf.Manifest
	if manifests, err = r.Manifests.Defaults(ctx); err != nil {
		l.Error(err, "Error reading YAML manifests")
		return ctrl.Result{}, err
	}
	git := module.Spec.Repository.Git
	if git == nil {
		manifests = manifests.Filter(mf.Not(mf.ByKind("GitRepository")))
	}
	if manifests, err = manifests.Transform(func(u *unstructured.Unstructured) error {
		u.SetName(module.GetName())
		u.SetNamespace(module.GetNamespace())
		if u.GetKind() == "GitRepository" {
			gitrepo := &sourcev1.GitRepository{}
			if err = r.Scheme.Convert(u, gitrepo, ctx); err != nil {
				return err
			}
			gitrepo.Spec.URL = git.URL
			gitrepo.Spec.Reference = &sourcev1.GitRepositoryRef{
				Branch: git.Branch,
				Tag:    git.Tag,
				Commit: git.Commit,
			}
			if err = controllerutil.SetOwnerReference(module, gitrepo, r.Scheme); err != nil {
				return err
			}
			if err = r.Scheme.Convert(gitrepo, u, nil); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		l.Error(err, "Error transforming a manifest resource")
		return ctrl.Result{}, err
	}
	if err = manifests.Apply(); err != nil {
		l.Error(err, "Error applying changes to resources in manifest")
		return ctrl.Result{}, err
	}
	l.Info("Successfully reconciled!")
	return ctrl.Result{}, nil
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
