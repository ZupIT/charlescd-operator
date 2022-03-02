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
	mf "github.com/manifestival/manifestival"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/ZupIT/charlescd-operator/internal/tracing"
)

const (
	sourceError              = "Error templating source"
	kustomizationValid       = "Kustomize configurations are valid"
	kustomizationSourceValid = "Kustomization were successfully rendered"
)

type (
	KustomizationClient interface {
		Kustomize(ctx context.Context, source, path string) (mf.Manifest, error)
	}
	KustomizationValidation struct {
		reconciler.Funcs

		kustomization KustomizationClient
		status        StatusWriter
	}
)

func NewKustomizationValidation(kustomization KustomizationClient, status StatusWriter) *KustomizationValidation {
	return &KustomizationValidation{kustomization: kustomization, status: status}
}

func (k *KustomizationValidation) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	module, ok := obj.(*charlescdv1alpha1.Module)
	if !ok || module.Spec.Kustomization == nil || !module.IsSourceReady() {
		return k.Next(ctx, obj)
	}

	return k.reconcile(ctx, module, module.Spec.Kustomization)
}

func (k *KustomizationValidation) reconcile(ctx context.Context, module *charlescdv1alpha1.Module, kustomization *charlescdv1alpha1.Kustomization) (ctrl.Result, error) {
	// check if this handler should act
	if module.Status.Source == nil || module.Status.Source.Path == "" {
		return k.Next(ctx, module)
	}

	// starting the context
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	log := logr.FromContextOrDiscard(ctx)

	source := module.Status.Source.Path
	path := kustomization.GitRepository.Path

	// templating Kustomization
	manifests, err := k.kustomization.Kustomize(ctx, source, path)
	if err != nil {
		log.Error(err, sourceError)
		if module.SetSourceInvalid(renderError, err.Error()) {
			return k.status.UpdateModuleStatus(ctx, module)
		}
		return k.Next(ctx, module)
	}

	// update status to success
	if module.SetSourceValid(kustomizationSourceValid) {
		return k.status.UpdateModuleStatus(ctx, module)
	}

	log.Info(kustomizationValid)
	return k.Next(contextWithResources(ctx, manifests.Resources()), module)
}
