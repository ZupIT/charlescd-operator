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
