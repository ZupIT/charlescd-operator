/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tracing

import (
	"context"
	"io"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
)

const (
	module  = "github.com/tiagoangelozup/charles-alpha"
	service = "charles"
)

func Initialize() (io.Closer, error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		return nil, err
	}
	provider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
		)),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
	)
	otel.SetTracerProvider(provider)
	return &closer{TracerProvider: provider}, nil
}

func SpanFromContext(ctx context.Context) Span {
	span := trace.SpanFromContext(ctx)
	return &defaultSpan{Span: span}
}

func StartSpanFromContext(ctx context.Context) (Span, context.Context) {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	spanName := strings.Replace(funcName, module+"/", "", 1)
	newCtx, span := otel.Tracer(service).Start(ctx, spanName)
	return &defaultSpan{Span: span}, newCtx
}
