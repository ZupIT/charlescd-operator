package tracing_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

func TestInitialize(t *testing.T) {
	t.Run("SHOULD be configurable WHEN using environment variables", func(t *testing.T) {
		_ = os.Setenv("JAEGER_SAMPLER_TYPE", "const")
		_ = os.Setenv("JAEGER_SAMPLER_PARAM", "1")
		closer, err := tracing.Initialize()
		if err != nil {
			t.Fatal(err)
		}
		defer closer.Close()

		span, _ := tracing.StartSpanFromContext(context.TODO())
		info := span.Info()
		assert.True(t, info.IsSampled)
	})
	t.Run("SHOULD be configurable WHEN using optional parameters", func(t *testing.T) {
		closer, err := tracing.Initialize(
			tracing.WithSamplerType(tracing.ConstantSampler),
			tracing.WithSamplerParam(1.0),
		)
		if err != nil {
			t.Fatal(err)
		}
		defer closer.Close()

		span, _ := tracing.StartSpanFromContext(context.TODO())
		info := span.Info()
		assert.True(t, info.IsSampled)
	})
}
