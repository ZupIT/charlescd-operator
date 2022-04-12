package runtime_test

import (
	climocks "github.com/ZupIT/charlescd-operator/internal/client/mocks"
	"github.com/ZupIT/charlescd-operator/internal/runtime"
	"github.com/ZupIT/charlescd-operator/internal/runtime/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Context("when getting manager client", func() {
		It("should return it successfully", func() {
			expectedClient := new(climocks.Client)
			mockManager := new(mocks.Manager)
			mockManager.On("GetClient").Return(expectedClient)
			client := runtime.Client(mockManager)
			Expect(client).To(Equal(expectedClient))
		})
	})
})
