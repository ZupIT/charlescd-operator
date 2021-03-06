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

package controllers

import (
	"github.com/google/wire"

	"github.com/ZupIT/charlescd-operator/internal/client"
	"github.com/ZupIT/charlescd-operator/internal/object"
	"github.com/ZupIT/charlescd-operator/internal/resources"
	"github.com/ZupIT/charlescd-operator/internal/runtime"
	"github.com/ZupIT/charlescd-operator/pkg/filter"
	"github.com/ZupIT/charlescd-operator/pkg/module"
	"github.com/ZupIT/charlescd-operator/pkg/transformer"
)

var providers = wire.NewSet( //nolint // used for compile-time dependency injection
	reconcilers,
	newModuleReconciler,
	client.Providers,
	filter.Providers,
	module.Providers,
	object.Providers,
	resources.Providers,
	runtime.Providers,
	transformer.Providers,
	wire.Bind(new(client.ManifestsReader), new(*resources.Manifests)),
	wire.Bind(new(module.GitRepositoryGetter), new(*client.GitRepository)),
	wire.Bind(new(module.HelmClient), new(*client.Helm)),
	wire.Bind(new(module.KustomizationClient), new(*client.Kustomization)),
	wire.Bind(new(module.ManifestClient), new(*client.Manifest)),
	wire.Bind(new(module.ManifestReader), new(*resources.Manifests)),
	wire.Bind(new(module.ObjectConverter), new(*object.UnstructuredConverter)),
	wire.Bind(new(module.StatusWriter), new(*client.Module)),
	wire.Bind(new(ModuleGetter), new(*client.Module)),
	wire.Bind(new(transformer.ObjectConverter), new(*object.UnstructuredConverter)),
	wire.Bind(new(transformer.ObjectReference), new(*object.Reference)),
	wire.Struct(new(ModuleHandler), "*"),
)

func reconcilers(m *ModuleReconciler) []Reconciler {
	return []Reconciler{m}
}
