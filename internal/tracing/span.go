package tracing

import (
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

type Span interface {
	Log(logr.Logger) logr.Logger
	Finish()
	Error(err error) error
	fmt.Stringer
}

type defaultSpan struct{ trace.Span }

func (s *defaultSpan) Log(logger logr.Logger) logr.Logger {
	ctx := s.SpanContext()
	if ctx.IsValid() {
		return logger.WithValues("trace", s.String())
	}
	return logger
}

func (s *defaultSpan) Finish() { s.End() }

func (s *defaultSpan) Error(err error) error {
	if err != nil {
		s.RecordError(err)
		s.SetStatus(codes.Error, err.Error())
		var serr *kerrors.StatusError
		if errors.As(err, &serr) {
			status := serr.Status()
			s.SetAttributes(attribute.Int64("code", int64(status.Code)))
			s.SetAttributes(attribute.String("reason", string(status.Reason)))
		}
	}
	return err
}

func (s *defaultSpan) String() string {
	ctx := s.SpanContext()
	if !ctx.IsValid() {
		return ""
	}
	traceID := ctx.TraceID()
	spanID := ctx.SpanID()
	isSampled := ctx.IsSampled()
	return fmt.Sprintf("%s:%s:%t", traceID, spanID, isSampled)
}
