package object

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type Reference struct {
	Scheme *runtime.Scheme
}

func (r *Reference) SetOwner(owner, object metav1.Object) error {
	if err := controllerutil.SetOwnerReference(owner, object, r.Scheme); err != nil {
		return fmt.Errorf("failed to set %T %q owner reference: %w", object, object.GetName(), err)
	}
	return nil
}

func (r *Reference) SetController(controller, object metav1.Object) error {
	if err := controllerutil.SetControllerReference(controller, object, r.Scheme); err != nil {
		return fmt.Errorf("failed to set %T %q controller reference: %w", object, object.GetName(), err)
	}
	return nil
}
