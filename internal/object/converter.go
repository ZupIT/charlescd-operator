package object

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type UnstructuredConverter struct {
	Scheme *runtime.Scheme
}

func (c *UnstructuredConverter) FromUnstructured(in *unstructured.Unstructured, out interface{}) error {
	if err := c.Scheme.Convert(in, out, context.TODO()); err != nil {
		return fmt.Errorf("error converting from unstructured object: %w", err)
	}
	return nil
}

func (c *UnstructuredConverter) ToUnstructured(in interface{}, out *unstructured.Unstructured) error {
	if err := c.Scheme.Convert(in, out, context.TODO()); err != nil {
		return fmt.Errorf("error converting to unstructured object: %w", err)
	}
	return nil
}
