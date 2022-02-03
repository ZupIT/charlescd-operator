package transformer

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type ObjectConverter interface {
	FromUnstructured(in *unstructured.Unstructured, out interface{}) error
	ToUnstructured(in interface{}, out *unstructured.Unstructured) error
}
