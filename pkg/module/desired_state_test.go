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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

var _ = Describe("DesiredState", func() {

	var ctx context.Context
	var manifestReaderMock *resources.Manifests
	var desiredState *module.DesiredState
	var clientMock *mocks.Client
	BeforeEach(func() {

		scheme := runtime.NewScheme()
		utilruntime.Must(clientgoscheme.AddToScheme(scheme))
		utilruntime.Must(charlescdv1alpha1.AddToScheme(scheme))
		utilruntime.Must(sourcev1beta1.AddToScheme(scheme))
		ctx = context.TODO()
		clientMock = new(mocks.Client)
		manifestReaderMock = &resources.Manifests{
			Client: clientMock,
		}
		gitRepository := &filter.GitRepository{}
		filters := &module.Filters{
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
		transformer := &module.Transformers{
			GitRepository: transformerGitRepository,
			Metadata:      metadata,
		}
		desiredState = module.NewDesiredState(filters, transformer, manifestReaderMock)
		reconciler.Chain(desiredState)
	})

	Context("when reconciling for the desired state", func() {
		It("should apply the correct desire state", func() {

			contextWithResources := fillContextResources(ctx)
			mod := newValidModule()

			clientMock.On("Get", mock.Anything).Return(nil, nil)
			clientMock.On("Create", mock.Anything, mock.Anything).Return(nil)
			_, err := desiredState.Reconcile(contextWithResources, mod)
			Expect(err).To(BeNil())
		})

		It("should apply the correct desire state", func() {
			applyError := errors.New("error applying manifests")
			contextWithResources := fillContextResources(ctx)
			mod := newValidModule()

			clientMock.On("Get", mock.Anything).Return(nil, nil)
			clientMock.On("Create", mock.Anything, mock.Anything).Return(applyError)
			_, err := desiredState.Reconcile(contextWithResources, mod)
			Expect(err).To(Equal(applyError))
		})

	})
})
