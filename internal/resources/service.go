package resources

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"

	"github.com/go-logr/logr"
	mf "github.com/manifestival/manifestival"

	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

//go:embed manifests.yaml
var manifests []byte

type Manifests struct{ Client mf.Client }

func (s *Manifests) LoadDefaults(ctx context.Context) (mf.Manifest, error) {
	span := tracing.SpanFromContext(ctx)
	log := logr.FromContextOrDiscard(ctx).V(1)
	reader := bytes.NewReader(manifests)
	m, err := mf.ManifestFrom(mf.Reader(reader), mf.UseClient(s.Client), mf.UseLogger(log))
	if err != nil {
		return mf.Manifest{}, span.Error(fmt.Errorf("failed to build deployments manifests: %w", err))
	}
	return m, nil
}
