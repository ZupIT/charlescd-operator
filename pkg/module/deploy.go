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
	"github.com/go-logr/logr"
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/ZupIT/charlescd-operator/internal/tracing"
	"github.com/angelokurtis/reconciler"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Deploy struct {
	reconciler.Funcs
	manifest ManifestReader
}

func NewDeploy(manifest ManifestReader) *Deploy {
	return &Deploy{manifest: manifest}
}

func (d *Deploy) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	module, ok := obj.(*charlescdv1alpha1.Module)
	if !ok {
		return d.Next(ctx, obj)
	}
	resources := resourcesFromContext(ctx)
	if len(resources) == 0 {
		return d.Next(ctx, obj)
	}
	return d.reconcile(ctx, module, resources)
}

func (d *Deploy) reconcile(ctx context.Context, module *charlescdv1alpha1.Module, resources []unstructured.Unstructured) (ctrl.Result, error) {
	manifests, err := d.manifests(ctx, resources)
	if err != nil {
		return d.RequeueOnErr(ctx, err)
	}

	manifests, err = d.transform(ctx, manifests, module)
	if err != nil {
		return d.RequeueOnErr(ctx, err)
	}

	err = d.apply(ctx, manifests)
	if err != nil {
		return d.RequeueOnErr(ctx, err)
	}

	return d.Next(ctx, module)
}

func (d *Deploy) apply(ctx context.Context, manifests mf.Manifest) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	err := manifests.Apply()
	if err != nil {
		return err
	}
	logr.FromContextOrDiscard(ctx).Info("Manifests applied")
	return err
}

func (d *Deploy) transform(ctx context.Context, manifests mf.Manifest, module *charlescdv1alpha1.Module) (mf.Manifest, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	return manifests.Transform(
		mf.InjectNamespace(module.GetNamespace()),
		mf.InjectOwner(module),
	)
}

func (d *Deploy) manifests(ctx context.Context, resources []unstructured.Unstructured) (mf.Manifest, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	return d.manifest.FromUnstructured(ctx, resources)
}
