package module

import (
	"context"

	"github.com/angelokurtis/reconciler"
	mf "github.com/manifestival/manifestival"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
	"github.com/tiagoangelozup/charles-alpha/pkg/filter"
	"github.com/tiagoangelozup/charles-alpha/pkg/transformer"
)

var logger = ctrl.Log.WithName("internal").WithName("module")

type (
	DesiredState struct {
		reconciler.Funcs

		filters      *Filters
		transformers *Transformers

		manifests Manifests
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

func NewDesiredState(filters *Filters, transformers *Transformers, manifests Manifests) *DesiredState {
	return &DesiredState{filters: filters, transformers: transformers, manifests: manifests}
}

func (d *DesiredState) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	if module, ok := obj.(*deployv1alpha1.Module); ok {
		return d.reconcile(ctx, module)
	}
	return d.Next(ctx, obj)
}

func (d *DesiredState) reconcile(ctx context.Context, module *deployv1alpha1.Module) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	l := logger.WithValues("trace", span)

	manifests, err := d.manifests.Defaults(ctx)
	if err != nil {
		l.Error(err, "Error reading YAML manifests")
		return d.RequeueOnErr(ctx, err)
	}

	// filters unnecessary manifests
	manifests = manifests.Filter(
		d.filters.FilterGitRepository(module),
	)

	// transform manifests to desired state
	if manifests, err = manifests.Transform(
		d.transformers.TransformMetadata(module),
		d.transformers.TransformGitRepository(module),
	); err != nil {
		l.Error(err, "Error transforming a manifest resource")
		return d.RequeueOnErr(ctx, err)
	}

	// apply desired state
	if err = manifests.Apply(); err != nil {
		l.Error(err, "Error applying changes to resources in manifest")
		return d.RequeueOnErr(ctx, err)
	}

	return d.Next(ctx, module)
}
