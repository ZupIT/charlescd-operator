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

package module

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// resourcesContextKey is how we find []unstructured.Unstructured in a context.Context.
type ResourcesContextKey struct{}

func contextWithResources(ctx context.Context, resources []unstructured.Unstructured) context.Context {
	return context.WithValue(ctx, ResourcesContextKey{}, resources)
}

func resourcesFromContext(ctx context.Context) []unstructured.Unstructured {
	if ctx == nil {
		return []unstructured.Unstructured{}
	}
	if v, ok := ctx.Value(ResourcesContextKey{}).([]unstructured.Unstructured); ok {
		return v
	}
	return []unstructured.Unstructured{}
}
