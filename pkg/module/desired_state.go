// Copyright 2022 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package module

import (
	"context"

	"github.com/ZupIT/charlescd-operator/internal/tracing"
	"github.com/angelokurtis/reconciler"
	"github.com/go-logr/logr"
	mf "github.com/manifestival/manifestival"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/ZupIT/charlescd-operator/pkg/filter"
	"github.com/ZupIT/charlescd-operator/pkg/transformer"
)

type (
	DesiredState struct {
		reconciler.Funcs

		filters      *Filters
		transformers *Transformers

		manifest ManifestReader
	}
	Transformers struct {
		*transformer.GitRepository
		*transformer.Metadata
	}
	Filters struct {
		*filter.GitRepository
	}
)

func NewDesiredState(filters *Filters, transformers *Transformers, manifest ManifestReader) *DesiredState {
	return &DesiredState{filters: filters, transformers: transformers, manifest: manifest}
}

func (d *DesiredState) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	if module, ok := obj.(*charlescdv1alpha1.Module); ok {
		return d.reconcile(ctx, module)
	}
	return d.Next(ctx, obj)
}

func (d *DesiredState) reconcile(ctx context.Context, module *charlescdv1alpha1.Module) (ctrl.Result, error) {
	manifests, err := d.manifests(ctx)
	if err != nil {
		return d.RequeueOnErr(ctx, err)
	}

	// filters unnecessary manifests
	manifests = d.filter(ctx, manifests, module)

	// transform manifests to desired state
	if manifests, err = d.transform(ctx, manifests, module); err != nil {
		return d.RequeueOnErr(ctx, err)
	}

	// apply desired state
	if err = d.apply(ctx, manifests); err != nil {
		return d.RequeueOnErr(ctx, err)
	}

	return d.Next(ctx, module)
}

func (d *DesiredState) apply(ctx context.Context, manifests mf.Manifest) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	err := manifests.Apply()
	if err != nil {
		return err
	}
	logr.FromContextOrDiscard(ctx).Info("Manifests applied")
	return err
}

func (d *DesiredState) transform(ctx context.Context, manifests mf.Manifest, module *charlescdv1alpha1.Module) (mf.Manifest, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	return manifests.Transform(
		d.transformers.TransformMetadata(module),
		d.transformers.TransformGitRepository(module),
	)
}

func (d *DesiredState) filter(ctx context.Context, manifests mf.Manifest, module *charlescdv1alpha1.Module) mf.Manifest {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	return manifests.Filter(
		d.filters.FilterGitRepository(module),
	)
}

func (d *DesiredState) manifests(ctx context.Context) (mf.Manifest, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	return d.manifest.LoadDefaults(ctx)
}
