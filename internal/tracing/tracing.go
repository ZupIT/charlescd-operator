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
	"os"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"go.opentelemetry.io/otel"
)

const (
	module      = "github.com/tiagoangelozup/charles-alpha"
	service     = "charles"
	envEndpoint = "OTEL_EXPORTER_JAEGER_ENDPOINT"
)

func Initialize() (io.Closer, error) {
	options := []tracesdk.TracerProviderOption{
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
		)),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
	}

	endpoint, ok := os.LookupEnv(envEndpoint)
	if ok {
		collectorEndpoint := jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint))
		e, err := jaeger.New(collectorEndpoint)
		if err != nil {
			return nil, err
		}
		options = append(options, tracesdk.WithBatcher(e))
	}

	tp := tracesdk.NewTracerProvider(options...)
	otel.SetTracerProvider(tp)
	return new(closer), nil
}

func SpanFromContext(ctx context.Context) Span {
	span := trace.SpanFromContext(ctx)
	return &defaultSpan{Span: span}
}

func StartSpanFromContext(ctx context.Context) (Span, context.Context) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	spanName := strings.Replace(funcName, module+"/", "", 1)
	newCtx, span := otel.Tracer(service).Start(ctx, spanName)
	s := &defaultSpan{Span: span}
	log := ctrl.Log.WithName(spanName).WithValues("trace", s.String())
	return s, logf.IntoContext(newCtx, log)
}
