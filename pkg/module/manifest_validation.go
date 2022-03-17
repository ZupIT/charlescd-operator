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
	"github.com/hashicorp/go-getter"
	mf "github.com/manifestival/manifestival"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/ZupIT/charlescd-operator/internal/tracing"
)

const manifestError = "ManifestLoadError"

type (
	ManifestClient interface {
		Load(ctx context.Context, releaseName, source, path string, values *apiextensionsv1.JSON) (mf.Manifest, error)
	}
	ManifestValidation struct {
		reconciler.Funcs
		status StatusWriter
	}
)

func NewManifestValidation(status StatusWriter) *ManifestValidation {
	return &ManifestValidation{status: status}
}

func (h *ManifestValidation) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	module, ok := obj.(*charlescdv1alpha1.Module)
	if !ok || module.Spec.Manifests == nil || !module.IsSourceReady() {
		return h.Next(ctx, obj)
	}
	return h.reconcile(ctx, module)
}

func (h *ManifestValidation) reconcile(ctx context.Context, module *charlescdv1alpha1.Module) (ctrl.Result, error) {
	// check if this handler should act
	if module.Status.Source == nil {
		return h.Next(ctx, module)
	}

	// starting the context
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()
	log := logr.FromContextOrDiscard(ctx)
	log.Info("Starting manifest validation reconcile")
	// Loading pure manifests
	dst, err := os.MkdirTemp(os.TempDir(), "manifests")

	err = getter.GetAny(dst, module.Status.Source.Path)
	if err != nil {
		log.Error(err, "Error loading manifests from source")
		if module.SetSourceInvalid(manifestError, err.Error()) {
			return h.status.UpdateModuleStatus(ctx, module)
		}
		return h.Next(ctx, module)
	}
	defer os.RemoveAll(dst)
	manifests, err := mf.NewManifest(dst)
	if err != nil {
		log.Error(err, "Error loading manifests from source")
		if module.SetSourceInvalid(manifestError, err.Error()) {
			return h.status.UpdateModuleStatus(ctx, module)
		}
		return h.Next(ctx, module)
	}

	// update status to success
	if module.SetSourceValid() {
		return h.status.UpdateModuleStatus(ctx, module)
	}

	log.WithValues("manifests", manifests.Resources()).Info("manifests is valid")
	return h.Next(contextWithResources(ctx, manifests.Resources()), module)
}
