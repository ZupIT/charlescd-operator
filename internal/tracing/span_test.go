package tracing_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

func TestSpan_Info(t *testing.T) {
	t.Parallel()

	t.Run("SHOULD show receiver and function name WHEN ask for operation name", func(t *testing.T) {
		t.Parallel()

		closer, err := tracing.Initialize()
		if err != nil {
			t.Fatal(err)
		}
		defer closer.Close()

		_, ctx := tracing.StartSpanFromContext(context.TODO())
		info := new(empty).f(ctx)

		assert.Equal(t, "tracing_test.(*empty).f", info.OperationName)
	})

	t.Run("SHOULD have a trace id WHEN child context is started", func(t *testing.T) {
		t.Parallel()

		closer, err := tracing.Initialize()
		if err != nil {
			t.Fatal(err)
		}
		defer closer.Close()

		_, ctx := tracing.StartSpanFromContext(context.TODO())
		info := new(empty).f(ctx)

		assert.Len(t, info.TraceID, 16)
	})

	t.Run("SHOULD have a span id WHEN child context is started", func(t *testing.T) {
		t.Parallel()

		closer, err := tracing.Initialize()
		if err != nil {
			t.Fatal(err)
		}
		defer closer.Close()

		_, ctx := tracing.StartSpanFromContext(context.TODO())
		info := new(empty).f(ctx)

		assert.Len(t, info.SpanID, 16)
	})

	t.Run("SHOULD generate a new span id WHEN child context is started", func(t *testing.T) {
		t.Parallel()

		closer, err := tracing.Initialize()
		if err != nil {
			t.Fatal(err)
		}
		defer closer.Close()

		_, ctx := tracing.StartSpanFromContext(context.TODO())
		info := new(empty).f(ctx)

		assert.NotEqual(t, info.SpanID, info.TraceID)
	})

	t.Run("SHOULD parent id must be equal the caller id WHEN child context is started", func(t *testing.T) {
		t.Parallel()

		closer, err := tracing.Initialize()
		if err != nil {
			t.Fatal(err)
		}
		defer closer.Close()

		_, ctx := tracing.StartSpanFromContext(context.TODO())
		info := new(empty).f(ctx)

		assert.Equal(t, info.ParentID, info.TraceID)
	})
}

type empty struct{}

func (e *empty) f(ctx context.Context) *tracing.Info {
	span, _ := tracing.StartSpanFromContext(ctx)
	info := span.Info()
	return info
}
