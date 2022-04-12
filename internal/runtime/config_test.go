package runtime_test

import (
	"github.com/ZupIT/charlescd-operator/internal/runtime/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"

	"github.com/ZupIT/charlescd-operator/internal/runtime"
)

var _ = Describe("Config", func() {
	Context("when getting config of client", func() {
		It("should return it successfully", func() {
			restConfig := new(rest.Config)
			mockManager := new(mocks.Manager)
			mockManager.On("GetConfig").Return(restConfig)
			config := runtime.Config(mockManager)
			Expect(config).To(Equal(restConfig))
		})
	})
})
