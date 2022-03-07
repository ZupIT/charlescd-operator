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

	"github.com/angelokurtis/reconciler"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	charlescdv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
	"github.com/tiagoangelozup/charles-alpha/pkg/filter"
	"github.com/tiagoangelozup/charles-alpha/pkg/transformer"
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
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := logr.FromContextOrDiscard(ctx)

	manifests, err := d.manifest.LoadDefaults(ctx)
	if err != nil {
		log.Error(err, "Error reading YAML manifests")
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
		log.Error(err, "Error transforming a manifest resource")
		return d.RequeueOnErr(ctx, err)
	}

	// apply desired state
	if err = manifests.Apply(); err != nil {
		log.Error(err, "Error applying changes to resources in manifest")
		return d.RequeueOnErr(ctx, err)
	}

	return d.Next(ctx, module)
}
