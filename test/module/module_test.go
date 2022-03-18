package module_test

import (
	"context"
	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/ZupIT/charlescd-operator/pkg/module"
	"github.com/ZupIT/charlescd-operator/test/module/mocks"
	"github.com/angelokurtis/reconciler"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("Module", func() {
	Context("when reconciling  pure manifests", func() {
		It("should not return error when are valid manifests", func() {
			ctx := context.TODO()
			statusWriterMock := new(mocks.StatusWriter)
			manifestClientMock := new(mocks.ManifestClient)
			manifestValidation := module.NewManifestValidation(statusWriterMock, manifestClientMock)
			reconciler.Chain(manifestValidation)
			module := setupModule()
			manifestClientMock.On("DownloadFromSource", mock.Anything, module.Status.Source.Path).Return("./data/deployment.yaml", nil)
			statusWriterMock.On("UpdateModuleStatus", mock.Anything, module).Return(ctrl.Result{}, nil)
			_, err := manifestValidation.Reconcile(ctx, module)
			Expect(err).To(BeNil())
			expectedCondition := metav1.Condition{
				Type:    charlescdv1alpha1.SourceReady,
				Status:  metav1.ConditionTrue,
				Reason:  "Validated",
				Message: "Helm chart templates were successfully rendered",
			}
			Expect(module.Status.Conditions).To(Equal([]metav1.Condition{expectedCondition}))
		})
	})

})

func setupModule() *charlescdv1alpha1.Module {
	module := new(charlescdv1alpha1.Module)
	module.Status.Conditions = []metav1.Condition{metav1.Condition{Type: "SourceReady", Status: metav1.ConditionTrue}}
	module.Spec.Manifests = &charlescdv1alpha1.Manifests{GitRepository: charlescdv1alpha1.GitRepository{URL: "https://example.com"}}
	module.Status.Source = &charlescdv1alpha1.Source{Path: "path/file.tgz"}
	return module
}
