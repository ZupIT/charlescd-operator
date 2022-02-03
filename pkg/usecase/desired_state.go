package usecase

import (
	"context"

	mf "github.com/manifestival/manifestival"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/runtime"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
	"github.com/tiagoangelozup/charles-alpha/pkg/filter"
	"github.com/tiagoangelozup/charles-alpha/pkg/transformer"
)

var logger = ctrl.Log.WithName("internal").WithName("usecase")

type (
	DesiredState struct {
		manifests Manifests

		transformers *Transformers
		filters      *Filters

		next runtime.ReconcilerOperation
	}
	Manifests interface {
		Defaults(ctx context.Context) (mf.Manifest, error)
	}
	Transformers struct {
		*transformer.GitRepository
		*transformer.Metadata
	}
	Filters struct {
		*filter.GitRepository
	}
)

func NewDesiredState(manifests Manifests, transformers *Transformers, filters *Filters) *DesiredState {
	return &DesiredState{manifests: manifests, transformers: transformers, filters: filters}
}

func (ds *DesiredState) SetNext(next runtime.ReconcilerOperation) {
	ds.next = next
}

func (ds *DesiredState) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	if module, ok := obj.(*deployv1alpha1.Module); ok {
		return ds.EnsureDesiredState(ctx, module)
	}
	return ds.next.Reconcile(ctx, obj)
}

func (ds *DesiredState) EnsureDesiredState(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	l := logger.WithValues("trace", span)

	manifests, err := ds.manifests.Defaults(ctx)
	if err != nil {
		l.Error(err, "Error reading YAML manifests")
		return runtime.RequeueOnErr(ctx, err)
	}

	// filters unnecessary manifests
	manifests = manifests.Filter(
		ds.filters.FilterGitRepository(module),
	)

	// transform manifests to desired state
	if manifests, err = manifests.Transform(
		ds.transformers.TransformMetadata(module),
		ds.transformers.TransformGitRepository(module),
	); err != nil {
		l.Error(err, "Error transforming a manifest resource")
		return runtime.RequeueOnErr(ctx, err)
	}

	// apply desired state
	if err = manifests.Apply(); err != nil {
		l.Error(err, "Error applying changes to resources in manifest")
		return runtime.RequeueOnErr(ctx, err)
	}

	return ds.next.Reconcile(ctx, module)
}
