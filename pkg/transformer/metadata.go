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

package transformer

import (
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
)

type Metadata struct{ reference ObjectReference }

func NewMetadata(reference ObjectReference) *Metadata {
	return &Metadata{reference: reference}
}

func (m *Metadata) TransformMetadata(module *charlescdv1alpha1.Module) mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		u.SetName(module.GetName())
		u.SetNamespace(module.GetNamespace())
		if err := m.reference.SetController(module, u); err != nil {
			return err
		}
		u.SetLabels(map[string]string{
			"app.kubernetes.io/managed-by": "charles",
			"app.kubernetes.io/name":       module.GetName(),
		})
		return nil
	}
}
