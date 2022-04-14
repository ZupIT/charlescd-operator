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

package module_test

import (
	"context"
	"errors"
	"github.com/ZupIT/charlescd-operator/pkg/module"
	"github.com/ZupIT/charlescd-operator/pkg/module/mocks"
	mf "github.com/manifestival/manifestival"

	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/angelokurtis/reconciler"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("Helm Validation Test", func() {
	var ctx context.Context
	var statusWriterMock *mocks.StatusWriter
	var helmClientMock *mocks.HelmClient
	var helmValidation *module.HelmValidation

	BeforeEach(func() {
		ctx = context.TODO()
		statusWriterMock = new(mocks.StatusWriter)
		helmClientMock = new(mocks.HelmClient)
		helmValidation = module.NewHelmValidation(helmClientMock, statusWriterMock)
		reconciler.Chain(helmValidation)
	})

	Context("when reconciling  helm charts", func() {
		It("should update status successfully when source are valid", func() {
			expectedCondition := metav1.Condition{
				Type:    charlescdv1alpha1.SourceValid,
				Status:  metav1.ConditionTrue,
				Reason:  "Validated",
				Message: "Helm chart templates were successfully rendered",
			}
			mod := setupHelmModule()

			helmClientMock.On(
				"Template",
				mock.Anything, mod.GetName(), mod.Status.Source.Path, mod.Spec.Helm.GitRepository.Path, mock.Anything,
			).Return(mf.Manifest{}, nil)
			statusWriterMock.On("UpdateModuleStatus", mock.Anything, mod).Return(ctrl.Result{}, nil)

			_, err := helmValidation.Reconcile(ctx, mod)

			Expect(err).To(BeNil())
			Expect(mod.Status.Conditions[1].Reason).To(Equal(expectedCondition.Reason))
			Expect(mod.Status.Conditions[1].Message).To(Equal(expectedCondition.Message))
			Expect(mod.Status.Conditions[1].Status).To(Equal(expectedCondition.Status))
			Expect(mod.Status.Conditions[1].Type).To(Equal(expectedCondition.Type))
		})

		It("should update status correctly when fails to download resource", func() {
			downloadError := errors.New("failed to download from source")
			mod := setupHelmModule()

			helmClientMock.On("Template",
				mock.Anything, mod.GetName(), mod.Status.Source.Path, mod.Spec.Helm.GitRepository.Path, mock.Anything,
			).Return(mf.Manifest{}, downloadError)

			statusWriterMock.On("UpdateModuleStatus", mock.Anything, mod).Return(ctrl.Result{}, nil)

			_, err := helmValidation.Reconcile(ctx, mod)

			Expect(err).To(BeNil())
			Expect(mod.Status.Conditions[1].Type).To(Equal(charlescdv1alpha1.SourceValid))
			Expect(mod.Status.Conditions[1].Status).To(Equal(metav1.ConditionFalse))
			Expect(mod.Status.Conditions[1].Message).To(Equal(downloadError.Error()))
		})
	})
})

func setupHelmModule() *charlescdv1alpha1.Module {
	module := new(charlescdv1alpha1.Module)
	module.Status.Conditions = []metav1.Condition{{Type: "SourceReady", Status: metav1.ConditionTrue}}
	module.Spec.Helm = &charlescdv1alpha1.Helm{GitRepository: &charlescdv1alpha1.GitRepository{URL: "https://example.com"}}
	module.Status.Source = &charlescdv1alpha1.Source{Path: "path/file.tgz"}
	return module
}
