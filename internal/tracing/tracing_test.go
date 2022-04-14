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

package tracing_test

import (
	"context"
	"errors"
	"github.com/ZupIT/charlescd-operator/internal/tracing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

var _ = Describe("Tracing", Ordered, func() {

	BeforeEach(func() {
		_, err := tracing.Initialize()
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when getting trace about context", func() {
		It("should return empty string if context has no trace id", func() {
			expectedString := ""
			span := tracing.SpanFromContext(context.TODO())
			spanString := span.String()
			Expect(spanString).To(Equal(expectedString))
			defer span.Finish()
		})

		It("should not return empty string if context has trace id", func() {
			span, _ := tracing.StartSpanFromContext(context.TODO())
			spanString := span.String()
			Expect(spanString).ToNot(BeEmpty())
			defer span.Finish()
		})

		It("should return empty string if context has no trace id", func() {
			spanError := errors.New("dummy error")
			span, _ := tracing.StartSpanFromContext(context.TODO())
			spanString := span.Error(spanError)
			Expect(spanString).To(Equal(spanError))
			defer span.Finish()
		})

		It("should return empty string if context has no trace id", func() {
			spanError := kerrors.NewBadRequest("bad request")
			span, _ := tracing.StartSpanFromContext(context.TODO())
			spanString := span.Error(spanError)
			Expect(spanString).To(Equal(spanError))
		})
	})
})
