package tracing

import (
	"runtime"
	"strings"

	"k8s.io/apimachinery/pkg/types"
)

type SpanOptions struct {
	operation string
	resource  *types.NamespacedName
}

func (o *SpanOptions) Operation() string {
	if o.operation == "" {
		pc, _, _, _ := runtime.Caller(2)
		funcName := runtime.FuncForPC(pc).Name()
		lastSlash := strings.LastIndexByte(funcName, '/')
		if lastSlash < 0 {
			lastSlash = 0
		} else {
			lastSlash++
		}
		return funcName[lastSlash:]
	}
	return o.operation
}

type SpanOptionFunc func(*SpanOptions)

func WithOperationName(operation string) SpanOptionFunc {
	return func(o *SpanOptions) {
		o.operation = operation
	}
}

func WithCustomResource(resource types.NamespacedName) SpanOptionFunc {
	return func(o *SpanOptions) {
		o.resource = &resource
	}
}
