package module

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// resourcesContextKey is how we find []unstructured.Unstructured in a context.Context.
type resourcesContextKey struct{}

func contextWithResources(ctx context.Context, resources []unstructured.Unstructured) context.Context {
	return context.WithValue(ctx, resourcesContextKey{}, resources)
}

func resourcesFromContext(ctx context.Context) []unstructured.Unstructured {
	if ctx == nil {
		return []unstructured.Unstructured{}
	}
	if v, ok := ctx.Value(resourcesContextKey{}).([]unstructured.Unstructured); ok {
		return v
	}
	return []unstructured.Unstructured{}
}
