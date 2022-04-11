package client_test

import (
	"context"
	"errors"
	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	moduleClient "github.com/ZupIT/charlescd-operator/internal/client"
	"github.com/ZupIT/charlescd-operator/internal/client/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Module", func() {
	var clientModule *moduleClient.Module
	var mockClient *mocks.Client
	var mockStatus *mocks.StatusWriter
	var module *charlescdv1alpha1.Module
	BeforeEach(func() {
		mockClient = new(mocks.Client)
		mockStatus = new(mocks.StatusWriter)
		clientModule = moduleClient.NewModule(mockClient)
		module = newModule()
	})

	Context("when updating the module status", func() {
		It("should not return error when module exists and is successfully updated", func() {
			mockClient.
				On("Get",
					mock.Anything,
					client.ObjectKey{Namespace: module.GetNamespace(), Name: module.GetName()},
					mock.Anything).
				Return(nil)
			mockClient.
				On("Status").
				Return(mockStatus)
			mockStatus.
				On("Patch",
					mock.Anything,
					mock.Anything,
					mock.Anything).
				Return(nil)
			_, err := clientModule.UpdateModuleStatus(context.TODO(), module)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error when the module fails to be updated", func() {
			mockClient.
				On("Get",
					mock.Anything,
					client.ObjectKey{Namespace: module.GetNamespace(), Name: module.GetName()},
					mock.Anything).
				Return(nil)
			mockClient.
				On("Status").
				Return(mockStatus)
			mockStatus.
				On("Patch",
					mock.Anything,
					mock.Anything,
					mock.Anything).
				Return(errors.New("error patching resource"))
			_, err := clientModule.UpdateModuleStatus(context.TODO(), module)
			Expect(err.Error()).To(Equal("failed to update Module status: error patching resource"))
		})

		It("should return error when fails to find the module", func() {
			mockClient.
				On("Get",
					mock.Anything,
					client.ObjectKey{Namespace: module.GetNamespace(), Name: module.GetName()},
					mock.Anything).
				Return(errors.New("error finding resource"))
			_, err := clientModule.UpdateModuleStatus(context.TODO(), module)
			Expect(err.Error()).To(Equal("failed to update Module status: failed to lookup resource: error finding resource"))
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

func newModule() *charlescdv1alpha1.Module {
	module := new(charlescdv1alpha1.Module)
	module.Status.Conditions = []metav1.Condition{}
	module.Spec.Manifests = &charlescdv1alpha1.Manifests{GitRepository: charlescdv1alpha1.GitRepository{URL: "https://example.com"}}
	module.Status.Source = &charlescdv1alpha1.Source{Path: "path/file.tgz"}
	return module
}
