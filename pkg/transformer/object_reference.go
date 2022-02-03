package transformer

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ObjectReference interface {
	SetOwner(owner, object metav1.Object) error
	SetController(controller, object metav1.Object) error
}
