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

package tracing

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

type Span interface {
	Finish()
	Error(err error) error
	fmt.Stringer
}

type defaultSpan struct{ trace.Span }

func (s *defaultSpan) setKubernetesResource(resource *kubernetesResource) {
	s.SetAttributes(attribute.String("kubernetes.resource.kind", resource.kind))
	s.SetAttributes(attribute.String("kubernetes.resource.name", resource.name))
	if resource.IsNamespaced() {
		s.SetAttributes(attribute.String("kubernetes.resource.namespace", resource.namespace))
	}
	s.SetAttributes(attribute.String("kubernetes.resource.version", resource.version))
}

func (s *defaultSpan) Finish() { s.End() }

func (s *defaultSpan) Error(err error) error {
	if err != nil {
		s.RecordError(err)
		s.SetStatus(codes.Error, err.Error())
		var serr *kerrors.StatusError
		if errors.As(err, &serr) {
			status := serr.Status()
			s.SetAttributes(attribute.Int64("kubernetes.error.code", int64(status.Code)))
			s.SetAttributes(attribute.String("kubernetes.error.reason", string(status.Reason)))
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
	return fmt.Sprintf("%s", traceID)
}
