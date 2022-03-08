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
