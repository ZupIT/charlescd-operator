package eventfilter

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
)

var moduleLog = ctrl.Log.WithName("eventfilter").
	WithValues("name", "modules", "version", "deploy.charlescd.io/v1alpha1").
	V(1)

type Module struct{}

func (m *Module) Create(event event.CreateEvent) bool {
	obj := event.Object
	if obj == nil {
		return false
	}
	if _, ok := obj.(*v1alpha1.Module); !ok {
		return false
	}
	moduleLog.Info("Resource created")
	return true
}

func (m *Module) Delete(event event.DeleteEvent) bool {
	obj := event.Object
	if obj == nil {
		return false
	}
	if _, ok := obj.(*v1alpha1.Module); !ok {
		return false
	}
	moduleLog.Info("Resource deleted")
	return true
}

func (m *Module) Update(event event.UpdateEvent) bool {
	if event.ObjectOld == nil || event.ObjectNew == nil {
		return false
	}
	if _, ok := event.ObjectOld.(*v1alpha1.Module); !ok {
		return false
	}
	if _, ok := event.ObjectNew.(*v1alpha1.Module); !ok {
		return false
	}
	moduleLog.Info("Resource updated")
	return true
}

func (m *Module) Generic(event event.GenericEvent) bool {
	obj := event.Object
	if obj == nil {
		return false
	}
	if _, ok := obj.(*v1alpha1.Module); !ok {
		return false
	}
	return true
}
