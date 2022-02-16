package module

import (
	"context"

	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ManifestReader interface {
	FromBytes(ctx context.Context, manifests []byte) (mf.Manifest, error)
	FromString(ctx context.Context, manifests string) (mf.Manifest, error)
	FromUnstructured(ctx context.Context, manifests []unstructured.Unstructured) (mf.Manifest, error)
	LoadDefaults(ctx context.Context) (mf.Manifest, error)
}
