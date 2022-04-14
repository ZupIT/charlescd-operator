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
