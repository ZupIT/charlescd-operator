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
