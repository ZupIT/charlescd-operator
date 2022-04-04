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

package resources

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/ZupIT/charlescd-operator/internal/tracing"
)

//go:embed manifests.yaml
var manifests []byte

type Manifests struct{ Client mf.Client }

func (s *Manifests) LoadDefaults(ctx context.Context) (mf.Manifest, error) {
	return s.FromBytes(ctx, manifests)
}

func (s *Manifests) FromBytes(ctx context.Context, manifests []byte) (mf.Manifest, error) {
	span := tracing.SpanFromContext(ctx)
	log := logr.FromContextOrDiscard(ctx).V(1)

	reader := bytes.NewReader(manifests)
	m, err := mf.ManifestFrom(mf.Reader(reader), mf.UseClient(s.Client), mf.UseLogger(log))
	if err != nil {
		return mf.Manifest{}, span.Error(fmt.Errorf("failed to read manifests from bytes: %w", err))
	}
	return m, nil
}

func (s *Manifests) FromString(ctx context.Context, manifests string) (mf.Manifest, error) {
	span := tracing.SpanFromContext(ctx)
	log := logr.FromContextOrDiscard(ctx).V(1)

	reader := strings.NewReader(manifests)
	m, err := mf.ManifestFrom(mf.Reader(reader), mf.UseClient(s.Client), mf.UseLogger(log))
	if err != nil {
		return mf.Manifest{}, span.Error(fmt.Errorf("failed to read manifests from string: %w", err))
	}
	return m, nil
}

func (s *Manifests) FromPath(ctx context.Context, path string, recursive bool) (mf.Manifest, error) {
	var source mf.Source
	if recursive {
		source = mf.Recursive(path)
	} else {
		source = mf.Path(path)
	}

	span := tracing.SpanFromContext(ctx)
	log := logr.FromContextOrDiscard(ctx).V(1)

	m, err := mf.ManifestFrom(source, mf.UseClient(s.Client), mf.UseLogger(log))
	if err != nil {
		return mf.Manifest{}, span.Error(fmt.Errorf("failed to read manifests from path: %w", err))
	}
	return m, nil
}

func (s *Manifests) FromUnstructured(ctx context.Context, manifests []unstructured.Unstructured) (mf.Manifest, error) {
	span := tracing.SpanFromContext(ctx)
	log := logr.FromContextOrDiscard(ctx).V(1)

	m, err := mf.ManifestFrom(mf.Slice(manifests), mf.UseClient(s.Client), mf.UseLogger(log))
	if err != nil {
		return mf.Manifest{}, span.Error(fmt.Errorf("failed to read manifests from string: %w", err))
	}
	return m, nil
}
