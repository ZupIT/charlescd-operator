package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

func StartSpanFromContext(ctx context.Context, options ...SpanOptionFunc) (*Span, context.Context) {
	opt := new(SpanOptions)

	for _, fn := range options {
		if fn == nil {
			continue
		}
		fn(opt)
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, opt.Operation())
	if opt.resource != nil {
		span.SetTag("kubernetes.resource", opt.resource.String())
	}
	return &Span{Span: span}, ctx
}

func SpanFromContext(ctx context.Context) *Span {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return nil
	}

	return &Span{Span: span}
}
