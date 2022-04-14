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

	"github.com/angelokurtis/reconciler"
	mf "github.com/manifestival/manifestival"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/ZupIT/charlescd-operator/internal/tracing"
)

const (
	renderError          = "RenderError"
	successRenderMessage = "Helm chart templates were successfully rendered"
)

type (
	HelmClient interface {
		Template(ctx context.Context, releaseName, source, path string, values *apiextensionsv1.JSON) (mf.Manifest, error)
	}
	HelmValidation struct {
		reconciler.Funcs

		helm   HelmClient
		status StatusWriter
	}
)

func NewHelmValidation(helm HelmClient, status StatusWriter) *HelmValidation {
	return &HelmValidation{helm: helm, status: status}
}

func (h *HelmValidation) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	module, ok := obj.(*charlescdv1alpha1.Module)
	if !ok || module.Spec.Helm == nil || !module.IsSourceReady() {
		return h.Next(ctx, obj)
	}
	return h.reconcile(ctx, module, module.Spec.Helm)
}

func (h *HelmValidation) reconcile(ctx context.Context, module *charlescdv1alpha1.Module, helm *charlescdv1alpha1.Helm) (ctrl.Result, error) {
	// check if this handler should act
	if module.Status.Source == nil || module.Status.Source.Path == "" {
		return h.Next(ctx, module)
	}

	source, path := module.Status.Source.Path, ""
	if helm.GitRepository != nil && helm.GitRepository.Path != "" {
		path = helm.GitRepository.Path
	}

	// templating Helm chart
	manifests, err := h.helmTemplate(ctx, module.GetName(), source, path, helm.Values)
	if err != nil {
		if module.SetSourceInvalid(renderError, err.Error()) {
			return h.status.UpdateModuleStatus(ctx, module)
		}
		return h.Next(ctx, module)
	}

	// update status to success
	if h.validate(ctx, module) {
		return h.status.UpdateModuleStatus(ctx, module)
	}

	return h.Next(contextWithResources(ctx, manifests.Resources()), module)
}

func (h *HelmValidation) validate(ctx context.Context, module *charlescdv1alpha1.Module) bool {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	logr.FromContextOrDiscard(ctx).Info("Helm chart is valid")
	return module.SetSourceValid(successRenderMessage)
}

func (h *HelmValidation) helmTemplate(ctx context.Context, name, source, path string, values *apiextensionsv1.JSON) (mf.Manifest, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	return h.helm.Template(ctx, name, source, path, values)
}
