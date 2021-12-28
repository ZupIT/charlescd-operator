package tracing

import (
	"errors"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

type (
	Span struct {
		opentracing.Span
	}
	Info struct {
		OperationName string
		TraceID       string
		SpanID        string
		ParentID      string
		IsSampled     bool
	}
)

func (s *Span) String() string {
	i := s.Info()
	return fmt.Sprintf("%s:%s:%s:%t", i.TraceID, i.SpanID, i.ParentID, i.IsSampled)
}

func (s *Span) Info() *Info {
	var operation string
	if js, ok := s.Span.(*jaeger.Span); ok {
		operation = js.OperationName()
	}
	if jsc, ok := s.Context().(jaeger.SpanContext); ok {
		return &Info{
			OperationName: operation,
			TraceID:       jsc.TraceID().String(),
			SpanID:        jsc.SpanID().String(),
			ParentID:      jsc.ParentID().String(),
			IsSampled:     jsc.IsSampled(),
		}
	}
	return nil
}

func (s *Span) SetError(err error) {
	ext.Error.Set(s, true)
	fields := make([]log.Field, 0)
	fields = append(fields, log.String("event", "error"), log.String("message", err.Error()))

	var serr *kerrors.StatusError
	if errors.As(err, &serr) {
		status := serr.Status()
		fields = append(fields, log.Int32("code", status.Code), log.String("reason", string(status.Reason)))
	}
	s.LogFields(fields...)
}

func (s *Span) HandleError(err error) error {
	if err != nil {
		s.SetError(err)
	}
	return err
}
