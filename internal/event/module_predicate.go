package event

import (
	"sigs.k8s.io/controller-runtime/pkg/event"

	"github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
)

type ModulePredicate struct{}

func NewModulePredicate() *ModulePredicate { return &ModulePredicate{} }

func (m *ModulePredicate) Create(event event.CreateEvent) bool {
	obj := event.Object
	if obj == nil {
		return false
	}
	if _, ok := obj.(*v1alpha1.Module); !ok {
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
	if _, ok := obj.(*v1alpha1.Module); !ok {
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
	moduleOld, ok := objOld.(*v1alpha1.Module)
	if !ok {
		return false
	}
	moduleNew, ok := objNew.(*v1alpha1.Module)
	if !ok {
		return false
	}
	log.WithValues("name", objNew.GetName(),
		"namespace", objNew.GetNamespace(),
		"resourceVersion", objNew.GetResourceVersion(),
		"diff", diff(
			&v1alpha1.Module{Spec: moduleOld.Spec, Status: moduleOld.Status},
			&v1alpha1.Module{Spec: moduleNew.Spec, Status: moduleNew.Status})).
		Info("Module updated")
	return true
}

func (m *ModulePredicate) Generic(event event.GenericEvent) bool {
	obj := event.Object
	if obj == nil {
		return false
	}
	if _, ok := obj.(*v1alpha1.Module); !ok {
		return false
	}
	return true
}
