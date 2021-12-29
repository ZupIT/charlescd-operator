package usecase

import (
	"context"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	mf "github.com/manifestival/manifestival"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrl "sigs.k8s.io/controller-runtime"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/runtime"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

var logger = ctrl.Log.WithName("internal").WithName("usecase")

type Manifests interface {
	Defaults(ctx context.Context) (mf.Manifest, error)
}

type ObjectConverter interface {
	FromUnstructured(in *unstructured.Unstructured, out interface{}) error
	ToUnstructured(in interface{}, out *unstructured.Unstructured) error
}

type ObjectReference interface {
	SetOwner(owner, object metav1.Object) error
	SetController(controller, object metav1.Object) error
}

type DesiredState struct {
	Manifests Manifests
	Object    ObjectConverter
	Reference ObjectReference
}

func (ds *DesiredState) EnsureDesiredState(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	l := logger.WithValues("trace", span)

	manifests, err := ds.Manifests.Defaults(ctx)
	if err != nil {
		l.Error(err, "Error reading YAML manifests")
		return runtime.RequeueOnErr(ctx, err)
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
			if err = ds.Object.FromUnstructured(u, gitrepo); err != nil {
				return err
			}
			gitrepo.Spec.URL = git.URL
			gitrepo.Spec.Reference = &sourcev1.GitRepositoryRef{
				Branch: git.Branch,
				Tag:    git.Tag,
				Commit: git.Commit,
			}
			if err = ds.Reference.SetController(module, gitrepo); err != nil {
				return err
			}
			if err = ds.Object.ToUnstructured(gitrepo, u); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		l.Error(err, "Error transforming a manifest resource")
		return runtime.RequeueOnErr(ctx, err)
	}
	if err = manifests.Apply(); err != nil {
		l.Error(err, "Error applying changes to resources in manifest")
		return runtime.RequeueOnErr(ctx, err)
	}
	return runtime.Finish()
}
