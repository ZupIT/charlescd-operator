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
	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/ZupIT/charlescd-operator/internal/object"
	"github.com/ZupIT/charlescd-operator/internal/resources"
	"github.com/ZupIT/charlescd-operator/pkg/filter"
	"github.com/ZupIT/charlescd-operator/pkg/module"
	"github.com/ZupIT/charlescd-operator/pkg/module/mocks"
	"github.com/ZupIT/charlescd-operator/pkg/transformer"
	"github.com/angelokurtis/reconciler"
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	mf "github.com/manifestival/manifestival"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var _ = Describe("DesiredState", func() {

	var ctx context.Context
	var manifestReaderMock *resources.Manifests
	var desiredState *module.DesiredState
	var clientMock *mocks.Client
	var filters *module.Filters
	var transformers *module.Transformers
	var scheme *runtime.Scheme
	BeforeEach(func() {

		scheme = runtime.NewScheme()
		utilruntime.Must(clientgoscheme.AddToScheme(scheme))
		utilruntime.Must(charlescdv1alpha1.AddToScheme(scheme))
		utilruntime.Must(sourcev1beta1.AddToScheme(scheme))
		ctx = context.TODO()
		clientMock = new(mocks.Client)
		manifestReaderMock = &resources.Manifests{
			Client: clientMock,
		}
		gitRepository := &filter.GitRepository{}
		filters = &module.Filters{
			GitRepository: gitRepository,
		}
		unstructuredConverter := &object.UnstructuredConverter{
			Scheme: scheme,
		}
		reference := &object.Reference{
			Scheme: scheme,
		}
		metadata := transformer.NewMetadata(reference)
		transformerGitRepository := transformer.NewGitRepository(unstructuredConverter)
		transformers = &module.Transformers{
			GitRepository: transformerGitRepository,
			Metadata:      metadata,
		}
		desiredState = module.NewDesiredState(filters, transformers, manifestReaderMock)
		reconciler.Chain(desiredState)
	})

	Context("when reconciling for the desired state", func() {
		It("should apply the correct desire state", func() {

			contextWithResources := fillContextResources(ctx)
			mod := newValidModule()

			clientMock.On("Get", mock.Anything).Return(nil, nil)
			clientMock.On("Create", mock.Anything, mock.Anything).Return(nil)
			_, err := desiredState.Reconcile(contextWithResources, mod)

			clientMock.AssertNumberOfCalls(GinkgoT(), "Get", 1)
			clientMock.AssertNumberOfCalls(GinkgoT(), "Create", 1)
			Expect(err).To(BeNil())
		})

		It("should return error when fails to load defaults manifests", func() {
			expectedError := errors.New("failed to load manifests")
			contextWithResources := fillContextResources(ctx)
			mod := newValidModule()
			manifestReaderMock := new(mocks.ManifestReader)
			desiredState := module.NewDesiredState(filters, transformers, manifestReaderMock)
			reconciler.Chain(desiredState)
			manifestReaderMock.On("LoadDefaults", mock.Anything).Return(mf.Manifest{}, expectedError)
			_, err := desiredState.Reconcile(contextWithResources, mod)

			clientMock.AssertNumberOfCalls(GinkgoT(), "Get", 0)
			clientMock.AssertNumberOfCalls(GinkgoT(), "Create", 0)
			Expect(err.Error()).To(Equal(expectedError.Error()))
		})

		It("should return error when fails to transform manifests", func() {
			manifests, _ := mf.ManifestFrom(mf.Slice(getGitRepositoryUnstructured()))
			expectedError := errors.New(
				`failed to set *unstructured.Unstructured "" controller reference: Object / is already owned by another GitRepository controller `)
			contextWithResources := fillContextResources(ctx)
			mod := newValidModule()
			manifestReaderMock := new(mocks.ManifestReader)
			desiredState := module.NewDesiredState(filters, transformers, manifestReaderMock)
			reconciler.Chain(desiredState)
			manifestReaderMock.On("LoadDefaults", mock.Anything).Return(manifests, nil)
			_, err := desiredState.Reconcile(contextWithResources, mod)
			clientMock.AssertNumberOfCalls(GinkgoT(), "Get", 0)
			clientMock.AssertNumberOfCalls(GinkgoT(), "Create", 0)
			Expect(err.Error()).To(Equal(expectedError.Error()))
		})

		It("should not do reconcile when resource is not a module", func() {

			contextWithResources := fillContextResources(ctx)
			mod := getGitRepository()
			clientMock.On("Get", mock.Anything).Return(nil, nil)
			clientMock.On("Create", mock.Anything, mock.Anything).Return(nil)
			_, err := desiredState.Reconcile(contextWithResources, mod)
			clientMock.AssertNumberOfCalls(GinkgoT(), "Get", 0)
			Expect(err).To(BeNil())
		})

		It("should return error when apply fails", func() {
			applyError := errors.New("error applying manifests")
			contextWithResources := fillContextResources(ctx)
			mod := newValidModule()

			clientMock.On("Get", mock.Anything).Return(nil, nil)
			clientMock.On("Create", mock.Anything, mock.Anything).Return(applyError)
			_, err := desiredState.Reconcile(contextWithResources, mod)
			Expect(err).To(Equal(applyError))

			clientMock.AssertNumberOfCalls(GinkgoT(), "Get", 1)
			clientMock.AssertNumberOfCalls(GinkgoT(), "Create", 1)
		})

	})
})

func getGitRepositoryUnstructured() []unstructured.Unstructured {
	unstructuredGit := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "source.toolkit.fluxcd.io/v1beta1",
			"kind":       "GitRepository",
			"metadata": map[string]interface{}{
				"name": "default",
			},
			"spec": map[string]interface{}{
				"interval": "1m",
				"url":      "",
			},
		},
	}
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(charlescdv1alpha1.AddToScheme(scheme))
	utilruntime.Must(sourcev1beta1.AddToScheme(scheme))
	error := controllerutil.SetControllerReference(getGitRepository(), &unstructuredGit, scheme)
	Expect(error).ToNot(HaveOccurred())
	return []unstructured.Unstructured{unstructuredGit}
}
