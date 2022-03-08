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
