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

	"github.com/angelokurtis/reconciler"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

<<<<<<< HEAD
	charlescdv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
=======
	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/ZupIT/charlescd-operator/internal/tracing"
>>>>>>> b9ad6cc8bbff9891be950e23f14133cd954d8f0b
)

type Status struct {
	reconciler.Funcs
	status StatusWriter
}

func NewStatus(status StatusWriter) *Status {
	return &Status{status: status}
}

func (s *Status) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	if module, ok := obj.(*charlescdv1alpha1.Module); ok {
		return s.reconcile(ctx, module)
	}
	return s.Next(ctx, obj)
}

func (s *Status) reconcile(ctx context.Context, module *charlescdv1alpha1.Module) (ctrl.Result, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	if module.UpdatePhase() {
		return s.status.UpdateModuleStatus(ctx, module)
	}
	return s.Next(ctx, module)
}
