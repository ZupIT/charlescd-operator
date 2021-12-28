package manifests

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"

	mf "github.com/manifestival/manifestival"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/tiagoangelozup/charles-alpha/internal/tracing"
)

var (
	//go:embed manifests.yaml
	manifests []byte
	logger    = ctrl.Log.WithName("manifests")
)

type Service struct{ Client mf.Client }

func (s *Service) Defaults(ctx context.Context) (mf.Manifest, error) {
	span := tracing.SpanFromContext(ctx)
	l := logger.V(1).WithValues("trace", span)
	reader := bytes.NewReader(manifests)
	m, err := mf.ManifestFrom(mf.Reader(reader), mf.UseClient(s.Client), mf.UseLogger(l))
	if err != nil {
		return mf.Manifest{}, span.HandleError(fmt.Errorf("failed to build deployments manifests: %w", err))
	}
	return m, nil
}
