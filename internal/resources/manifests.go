package resources

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	mf "github.com/manifestival/manifestival"

	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
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
