package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter"
	mf "github.com/manifestival/manifestival"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
)

type (
	ManifestsReader interface {
		FromString(ctx context.Context, manifests string) (mf.Manifest, error)
	}
	Helm struct{ manifests ManifestsReader }
)

func NewHelm(manifests ManifestsReader) *Helm {
	return &Helm{manifests: manifests}
}

func (h *Helm) Template(ctx context.Context, releaseName, source, path string) (mf.Manifest, error) {
	destination := source[0 : len(source)-len(".tar.gz")]
	if path != "" {
		source += "//" + path
		destination = filepath.Join(destination, path)
	}

	if err := getter.GetAny(destination, source); err != nil {
		return mf.Manifest{}, fmt.Errorf("error extracting Source artifact %s: %w", source, err)
	}
	defer os.RemoveAll(destination)

	r, err := h.template(ctx, releaseName, destination, map[string]interface{}{})
	if err != nil {
		return mf.Manifest{}, fmt.Errorf("error rendering Helm chart templates: %w", err)
	}

	return h.manifests.FromString(ctx, r.Manifest)
}

func (h *Helm) template(ctx context.Context, releaseName, chart string, values map[string]interface{}) (*release.Release, error) {
	c, err := loader.LoadDir(chart)
	if err != nil {
		return nil, fmt.Errorf("error rendering Helm chart templates: %w", err)
	}

	config := new(action.Configuration)
	client := action.NewInstall(config)
	client.DryRun = true
	client.ReleaseName = releaseName
	client.Replace = true
	client.ClientOnly = true
	client.IncludeCRDs = true

	return client.RunWithContext(ctx, c, values)
}
