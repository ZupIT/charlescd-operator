package module_test

import (
	"context"
	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/ZupIT/charlescd-operator/pkg/module/mocks"
	"github.com/angelokurtis/reconciler"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/ZupIT/charlescd-operator/pkg/module"
)

var _ = Describe("Status", func() {
	var ctx context.Context
	var statusWriterMock *mocks.StatusWriter
	var status *module.Status
	BeforeEach(func() {

		ctx = context.TODO()
		statusWriterMock = new(mocks.StatusWriter)
		status = module.NewStatus(statusWriterMock)
		reconciler.Chain(status)
	})
	Context("when reconciling  pure manifests", func() {
		It("should update phase successfully when source are ready", func() {
			exoectedStatus := "Ready"
			mod := newValidModule()
			statusWriterMock.On("UpdateModuleStatus", mock.Anything, mod).Return(ctrl.Result{}, nil)
			_, err := status.Reconcile(ctx, mod)

			Expect(err).To(BeNil())
			Expect(mod.Status.Phase).To(Equal(exoectedStatus))

		})

		It("should update phase successfully when source are ready", func() {
			exoectedStatus := "Failed"
			mod := newNotValidModule()
			statusWriterMock.On("UpdateModuleStatus", mock.Anything, mod).Return(ctrl.Result{}, nil)
			_, err := status.Reconcile(ctx, mod)

			Expect(err).To(BeNil())
			Expect(mod.Status.Phase).To(Equal(exoectedStatus))

		})

	})
})

func newNotValidModule() *charlescdv1alpha1.Module {
	module := new(charlescdv1alpha1.Module)
	module.Status.Conditions = []metav1.Condition{{Type: "SourceReady", Status: metav1.ConditionFalse}, {Type: "SourceValid", Status: metav1.ConditionFalse}}
	module.Spec.Manifests = &charlescdv1alpha1.Manifests{GitRepository: charlescdv1alpha1.GitRepository{URL: "https://example.com"}}
	module.Status.Source = &charlescdv1alpha1.Source{Path: "path/file.tgz"}
	return module
}
