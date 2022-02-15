package client

import (
	"context"

	mf "github.com/manifestival/manifestival"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"k8s.io/apimachinery/pkg/runtime"
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

func (h *Helm) Template(ctx context.Context, name, chart string, values runtime.RawExtension) (mf.Manifest, error) {
	chartRequested, err := loader.LoadDir(chart)
	if err != nil {
		return mf.Manifest{}, err
	}

	act := templateAction(name)
	release, err := act.RunWithContext(ctx, chartRequested, map[string]interface{}{})
	if err != nil || release == nil {
		return mf.Manifest{}, err
	}

	return h.manifests.FromString(ctx, release.Manifest)
}

func templateAction(name string) *action.Install {
	config := new(action.Configuration)
	client := action.NewInstall(config)
	client.DryRun = true
	client.ReleaseName = name
	client.Replace = true
	client.ClientOnly = true
	client.IncludeCRDs = true
	return client
}
