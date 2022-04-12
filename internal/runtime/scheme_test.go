package runtime_test

import (
	"github.com/ZupIT/charlescd-operator/internal/runtime"
	"github.com/ZupIT/charlescd-operator/internal/runtime/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("Scheme", func() {
	Context("when getting scheme of client", func() {
		It("should return it successfully", func() {
			expectedScheme := k8sruntime.NewScheme()
			mockManager := new(mocks.Manager)
			mockManager.On("GetScheme").Return(expectedScheme)
			scheme := runtime.Scheme(mockManager)
			Expect(scheme).To(Equal(expectedScheme))
		})
	})
})
