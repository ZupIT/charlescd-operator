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
	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/ZupIT/charlescd-operator/internal/object"
	"github.com/ZupIT/charlescd-operator/internal/resources"
	"github.com/ZupIT/charlescd-operator/pkg/module"
	"github.com/ZupIT/charlescd-operator/pkg/module/mocks"
	"github.com/angelokurtis/reconciler"
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	mf "github.com/manifestival/manifestival"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

type resourcesContextKey struct{}

var _ = Describe("CheckComponents", func() {
	var ctx context.Context
	var statusWriterMock *mocks.StatusWriter
	var checkComponents *module.CheckComponents

	BeforeEach(func() {
		scheme := runtime.NewScheme()
		utilruntime.Must(clientgoscheme.AddToScheme(scheme))
		utilruntime.Must(charlescdv1alpha1.AddToScheme(scheme))
		utilruntime.Must(sourcev1beta1.AddToScheme(scheme))
		ctx = context.TODO()
		statusWriterMock = new(mocks.StatusWriter)
		unstructuredConverter := &object.UnstructuredConverter{
			Scheme: scheme,
		}
		manifests := &resources.Manifests{}
		checkComponents = module.NewCheckComponents(manifests, unstructuredConverter, statusWriterMock)
		reconciler.Chain(checkComponents)
	})

	Context("when reconciling for checking components", func() {
		It("should update status successfully with deployable components", func() {

			contextWithResources := fillContextResources(ctx)
			mod := newValidModule()

			statusWriterMock.On("UpdateModuleStatus", mock.Anything, mod).Return(ctrl.Result{}, nil)
			_, err := checkComponents.Reconcile(contextWithResources, mod)
			Expect(err).To(BeNil())
			Expect(mod.Status.Components[0].Name).To(Equal("quiz-app"))
			Expect(mod.Status.Components[0].Containers[0].Name).To(Equal("quiz-app"))
			Expect(mod.Status.Components[0].Containers[0].Image).To(Equal("charlescd/quiz-app:1.0"))
		})

		It("should not update status when are no resources in context", func() {

			contextWIthoutResources := context.TODO()
			mod := newValidModule()

			statusWriterMock.On("UpdateModuleStatus", mock.Anything, mod).Return(ctrl.Result{}, nil)
			_, err := checkComponents.Reconcile(contextWIthoutResources, mod)
			Expect(err).To(BeNil())
			statusWriterMock.AssertNumberOfCalls(GinkgoT(), "UpdateModuleStatus", 0)
		})

		It("should not update status when source is not ready", func() {

			contextWithResources := fillContextResources(ctx)
			mod := newNotReadyModule()

			statusWriterMock.On("UpdateModuleStatus", mock.Anything, mod).Return(ctrl.Result{}, nil)
			_, err := checkComponents.Reconcile(contextWithResources, mod)
			Expect(err).To(BeNil())
			statusWriterMock.AssertNumberOfCalls(GinkgoT(), "UpdateModuleStatus", 0)
		})

		It("should update status when have duplicated components name", func() {
			expectedCondition := metav1.Condition{
				Type:    charlescdv1alpha1.SourceValid,
				Status:  metav1.ConditionFalse,
				Reason:  "DuplicatedComponent",
				Message: "test-deploy: component already present on module",
			}
			contextWithResources := fillContextWithDuplicatedResources(ctx)
			mod := newValidModule()

			statusWriterMock.On("UpdateModuleStatus", mock.Anything, mod).Return(ctrl.Result{}, nil)
			_, err := checkComponents.Reconcile(contextWithResources, mod)
			Expect(err).To(BeNil())
			Expect(mod.Status.Conditions[1].Status).To(Equal(expectedCondition.Status))
			Expect(mod.Status.Conditions[1].Message).To(Equal(expectedCondition.Message))
			Expect(mod.Status.Conditions[1].Reason).To(Equal(expectedCondition.Reason))
			Expect(mod.Status.Conditions[1].Type).To(Equal(expectedCondition.Type))
		})
	})
})

func fillContextResources(ctx context.Context) context.Context {
	manifests, err := mf.ManifestFrom(mf.Reader(strings.NewReader(getArtifactData())))
	Expect(err).To(BeNil())
	return context.WithValue(ctx, module.ResourcesContextKey{}, manifests.Resources())
}

func fillContextWithDuplicatedResources(ctx context.Context) context.Context {
	resourceDeployment := getUnstructuredDeployment()

	resourceDeploymentDuplicated := getUnstructuredDeployment()
	return context.WithValue(ctx, module.ResourcesContextKey{}, []unstructured.Unstructured{resourceDeployment, resourceDeploymentDuplicated})
}

func getUnstructuredDeployment() unstructured.Unstructured {
	u := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":      "test-deploy",
				"namespace": "test",
			},
		},
	}
	container := map[string]interface{}{
		"name":  "web",
		"image": "nginx:1.12",
	}
	containers := make([]interface{}, 1)
	containers[0] = container

	err := unstructured.SetNestedField(u.Object, containers, "spec", "template", "spec", "containers")
	Expect(err).ToNot(HaveOccurred())
	return u
}

func newValidModule() *charlescdv1alpha1.Module {
	module := new(charlescdv1alpha1.Module)
	module.Status.Conditions = []metav1.Condition{{Type: "SourceReady", Status: metav1.ConditionTrue}, {Type: "SourceValid", Status: metav1.ConditionTrue}}
	module.Spec.Manifests = &charlescdv1alpha1.Manifests{GitRepository: charlescdv1alpha1.GitRepository{URL: "https://example.com"}}
	module.Status.Source = &charlescdv1alpha1.Source{Path: "path/file.tgz"}
	return module
}
func newNotReadyModule() *charlescdv1alpha1.Module {
	module := new(charlescdv1alpha1.Module)
	module.Spec.Manifests = &charlescdv1alpha1.Manifests{GitRepository: charlescdv1alpha1.GitRepository{URL: "https://example.com"}}
	module.Status.Source = &charlescdv1alpha1.Source{Path: "path/file.tgz"}
	return module
}
