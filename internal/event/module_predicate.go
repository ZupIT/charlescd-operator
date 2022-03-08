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

package event

import (
	"sigs.k8s.io/controller-runtime/pkg/event"

	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
)

type ModulePredicate struct{}

func NewModulePredicate() *ModulePredicate { return &ModulePredicate{} }

func (m *ModulePredicate) Create(event event.CreateEvent) bool {
	obj := event.Object
	if obj == nil {
		return false
	}
	if _, ok := obj.(*charlescdv1alpha1.Module); !ok {
		return false
	}
	log.WithValues("name", obj.GetName(),
		"namespace", obj.GetNamespace(),
		"resourceVersion", obj.GetResourceVersion()).
		Info("Module created")
	return true
}

func (m *ModulePredicate) Delete(event event.DeleteEvent) bool {
	obj := event.Object
	if obj == nil {
		return false
	}
	if _, ok := obj.(*charlescdv1alpha1.Module); !ok {
		return false
	}
	log.WithValues("name", obj.GetName(),
		"namespace", obj.GetNamespace(),
		"resourceVersion", obj.GetResourceVersion()).
		Info("Module deleted")
	return true
}

func (m *ModulePredicate) Update(event event.UpdateEvent) bool {
	objOld, objNew := event.ObjectOld, event.ObjectNew
	if objOld == nil || objNew == nil {
		return false
	}
	moduleOld, ok := objOld.(*charlescdv1alpha1.Module)
	if !ok {
		return false
	}
	moduleNew, ok := objNew.(*charlescdv1alpha1.Module)
	if !ok {
		return false
	}
	log.WithValues("name", objNew.GetName(),
		"namespace", objNew.GetNamespace(),
		"resourceVersion", objNew.GetResourceVersion(),
		"diff", diff(
			&charlescdv1alpha1.Module{Spec: moduleOld.Spec, Status: moduleOld.Status},
			&charlescdv1alpha1.Module{Spec: moduleNew.Spec, Status: moduleNew.Status})).
		Info("Module updated")
	return true
}

func (m *ModulePredicate) Generic(event event.GenericEvent) bool {
	obj := event.Object
	if obj == nil {
		return false
	}
	if _, ok := obj.(*charlescdv1alpha1.Module); !ok {
		return false
	}
	return true
}
