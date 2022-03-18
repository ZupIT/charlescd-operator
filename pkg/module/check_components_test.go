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
	"gopkg.in/h2non/gock.v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			gock.New("https://example.com").
				Get("/manifests").
				Reply(200).
				BodyString(getArtifactData())

			contextWithResources := fillContextResources(ctx)
			mod := newValidModule()

			statusWriterMock.On("UpdateModuleStatus", mock.Anything, mod).Return(ctrl.Result{}, nil)
			_, err := checkComponents.Reconcile(contextWithResources, mod)
			Expect(err).To(BeNil())
			Expect(mod.Status.Components[0].Name).To(Equal("quiz-app"))
			Expect(mod.Status.Components[0].Containers[0].Name).To(Equal("quiz-app"))
			Expect(mod.Status.Components[0].Containers[0].Image).To(Equal("charlescd/quiz-app:1.0"))
		})
	})
})

func fillContextResources(ctx context.Context) context.Context {
	manifests, err := mf.ManifestFrom(mf.Reader(strings.NewReader(getArtifactData())))
	Expect(err).To(BeNil())
	return context.WithValue(ctx, module.ResourcesContextKey{}, manifests.Resources())
}

func newValidModule() *charlescdv1alpha1.Module {
	module := new(charlescdv1alpha1.Module)
	module.Status.Conditions = []metav1.Condition{{Type: "SourceReady", Status: metav1.ConditionTrue}, {Type: "SourceValid", Status: metav1.ConditionTrue}}
	module.Spec.Manifests = &charlescdv1alpha1.Manifests{GitRepository: charlescdv1alpha1.GitRepository{URL: "https://example.com"}}
	module.Status.Source = &charlescdv1alpha1.Source{Path: "path/file.tgz"}
	return module
}
