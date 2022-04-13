package tracing

import (
	"context"
	"fmt"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type kubernetesResourceContextKey struct{}

type kubernetesResource struct {
	kind, name, namespace, version string
}

func (k *kubernetesResource) IsNamespaced() bool {
	return len(k.namespace) > 0
}

func (k *kubernetesResource) String() string {
	if k.IsNamespaced() {
		return fmt.Sprintf("%s/%s/%s/%s", k.kind, k.namespace, k.name, k.version)
	}
	return fmt.Sprintf("%s/%s/%s", k.kind, k.name, k.version)
}

type ContextOptionFunc func(context.Context) context.Context

func WithResource(obj client.Object) ContextOptionFunc {
	return func(ctx context.Context) context.Context {
		gvk := obj.GetObjectKind().GroupVersionKind()
		kind := fmt.Sprintf("%s.%s", strings.ToLower(gvk.Kind), gvk.Group)
		if len(gvk.Group) == 0 {
			kind = strings.ToLower(gvk.Kind)
		}
		return context.WithValue(ctx, kubernetesResourceContextKey{}, &kubernetesResource{
			kind:      kind,
			name:      obj.GetName(),
			namespace: obj.GetNamespace(),
			version:   obj.GetResourceVersion(),
		})
	}
}
