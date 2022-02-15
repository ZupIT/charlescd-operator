package module

import (
	"context"

	mf "github.com/manifestival/manifestival"
)

type ManifestReader interface {
	FromString(ctx context.Context, manifests string) (mf.Manifest, error)
	LoadDefaults(ctx context.Context) (mf.Manifest, error)
}
