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

package object_test

import (
	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/ZupIT/charlescd-operator/internal/object"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"
)

var _ = Describe("Reference", func() {
	var objectReference *object.Reference
	BeforeEach(func() {
		scheme := runtime.NewScheme()
		utilruntime.Must(clientgoscheme.AddToScheme(scheme))
		utilruntime.Must(charlescdv1alpha1.AddToScheme(scheme))
		objectReference = &object.Reference{Scheme: scheme}
	})
	Context("when putting a owner for a resource", func() {
		It("should do it sucessfully when the owner is valid", func() {
			object := getUnstructuredDeployment()
			owner := newValidModule()
			err := objectReference.SetOwner(owner, object)
			Expect(err).ToNot(HaveOccurred())
			references := object.GetOwnerReferences()
			for _, ref := range references {
				Expect(ref.Name).To(Equal(owner.GetName()))
			}
		})

		It("should return error when is a cross namespace owner reference", func() {
			crossNamespaceError := `failed to set *unstructured.Unstructured "test-deploy" ` +
				`owner reference: cross-namespace owner references are disallowed, owner's namespace some-namespace, obj's namespace test`
			object := getUnstructuredDeployment()
			owner := newModuleWithDifferentNamespace()
			err := objectReference.SetOwner(owner, object)
			Expect(err.Error()).To(Equal(crossNamespaceError))
		})
	})

	Context("when putting a controller for a resource", func() {
		It("should return error when is a cross namespace owner reference", func() {
			expectedBoolean := pointer.Bool(true)
			object := getUnstructuredDeployment()
			owner := newValidModule()
			err := objectReference.SetController(owner, object)
			Expect(err).ToNot(HaveOccurred())
			references := object.GetOwnerReferences()
			for _, ref := range references {
				Expect(ref.BlockOwnerDeletion).To(Equal(expectedBoolean))
				Expect(ref.Controller).To(Equal(expectedBoolean))
				Expect(ref.Name).To(Equal(owner.GetName()))
			}
		})
		It("should return error when is a cross namespace owner reference", func() {
			crossNamespaceError := `failed to set *unstructured.Unstructured "test-deploy" ` +
				`controller reference: cross-namespace owner references are disallowed, owner's namespace some-namespace, obj's namespace test`
			object := getUnstructuredDeployment()
			owner := newModuleWithDifferentNamespace()
			err := objectReference.SetController(owner, object)
			Expect(err.Error()).To(Equal(crossNamespaceError))
		})
	})
})

func newValidModule() *charlescdv1alpha1.Module {
	module := new(charlescdv1alpha1.Module)
	module.Status.Conditions = []metav1.Condition{{Type: "SourceReady", Status: metav1.ConditionTrue}, {Type: "SourceValid", Status: metav1.ConditionTrue}}
	module.Spec.Manifests = &charlescdv1alpha1.Manifests{GitRepository: charlescdv1alpha1.GitRepository{URL: "https://example.com"}}
	module.Status.Source = &charlescdv1alpha1.Source{Path: "path/file.tgz"}
	module.SetNamespace("test")
	module.SetName("test-module")
	return module
}

func newModuleWithDifferentNamespace() *charlescdv1alpha1.Module {
	module := new(charlescdv1alpha1.Module)
	module.Status.Conditions = []metav1.Condition{{Type: "SourceReady", Status: metav1.ConditionTrue}, {Type: "SourceValid", Status: metav1.ConditionTrue}}
	module.Spec.Manifests = &charlescdv1alpha1.Manifests{GitRepository: charlescdv1alpha1.GitRepository{URL: "https://example.com"}}
	module.Status.Source = &charlescdv1alpha1.Source{Path: "path/file.tgz"}
	module.SetNamespace("some-namespace")
	module.SetName("test-module")
	return module
}
