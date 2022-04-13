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
