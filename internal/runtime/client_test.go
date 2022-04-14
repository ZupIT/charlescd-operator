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
