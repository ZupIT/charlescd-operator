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
	"errors"
	"fmt"

	"github.com/angelokurtis/reconciler"
	"github.com/go-logr/logr"
	mf "github.com/manifestival/manifestival"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/ZupIT/charlescd-operator/internal/tracing"
)

var ErrorDuplicatedComponent = errors.New("component already present on module")

const (
	DuplicatedComponent = "DuplicatedComponent"
)

type (
	ObjectConverter interface {
		FromUnstructured(in *unstructured.Unstructured, out interface{}) error
	}
	CheckComponents struct {
		reconciler.Funcs
		manifest ManifestReader
		object   ObjectConverter
		status   StatusWriter
	}
)

func NewCheckComponents(manifest ManifestReader, object ObjectConverter, status StatusWriter) *CheckComponents {
	return &CheckComponents{manifest: manifest, object: object, status: status}
}

func (c *CheckComponents) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	module, ok := obj.(*charlescdv1alpha1.Module)
	if !ok || !module.IsSourceValid() || !module.IsSourceReady() {
		return c.Next(ctx, obj)
	}
	resources := resourcesFromContext(ctx)
	if len(resources) == 0 {
		return c.Next(ctx, obj)
	}
	return c.reconcile(ctx, module, resources)
}

func (c *CheckComponents) reconcile(ctx context.Context, module *charlescdv1alpha1.Module, resources []unstructured.Unstructured) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := logr.FromContextOrDiscard(ctx)

	manifests, err := c.manifest.FromUnstructured(ctx, resources)
	if err != nil {
		return c.RequeueOnErr(ctx, err)
	}

	components := make([]*charlescdv1alpha1.Component, 0, 0)
	for _, u := range manifests.Filter(mf.ByKind("Deployment")).Resources() {
		deploy := &appsv1.Deployment{}
		if err := c.object.FromUnstructured(&u, deploy); err != nil {
			return c.RequeueOnErr(ctx, err)
		}
		component := &charlescdv1alpha1.Component{Name: deploy.GetName()}
		for _, container := range deploy.Spec.Template.Spec.Containers {
			component.Containers = append(component.Containers, &charlescdv1alpha1.Container{
				Name:  container.Name,
				Image: container.Image,
			})
			if err := c.checkComponentIsAlreadyPresent(components, component); err != nil {
				module.SetSourceInvalid(DuplicatedComponent, err.Error())
				return c.status.UpdateModuleStatus(ctx, module)
			}
			components = append(components, component)
		}
	}

	if total := len(components); total > 0 {
		log.Info("Deployable components found", "total", total)
	} else {
		log.Info("No deployable components were found")
	}

	if module.SetComponents(components) {
		return c.status.UpdateModuleStatus(ctx, module)
	}
	return c.Next(ctx, module)
}

func (c *CheckComponents) checkComponentIsAlreadyPresent(components []*charlescdv1alpha1.Component, component *charlescdv1alpha1.Component) error {
	for _, c := range components {
		if c.Name == component.Name {
			return fmt.Errorf("%s: %w", c.Name, ErrorDuplicatedComponent)
		}
	}
	return nil
}
