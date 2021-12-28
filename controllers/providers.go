/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"github.com/google/wire"

	"github.com/tiagoangelozup/charles-alpha/internal/manifests"
	"github.com/tiagoangelozup/charles-alpha/internal/module"
	"github.com/tiagoangelozup/charles-alpha/internal/runtime"
)

var providers = wire.NewSet(
	reconcilers,
	manifests.Providers,
	module.Providers,
	runtime.Providers,
	wire.Bind(new(Manifests), new(*manifests.Service)),
	wire.Bind(new(ModuleGetter), new(*module.Service)),
	wire.Struct(new(ModuleReconciler), "*"),
)

func reconcilers(m *ModuleReconciler) []Reconciler {
	return []Reconciler{m}
}