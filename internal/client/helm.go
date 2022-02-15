package client

import (
	"context"

	mf "github.com/manifestival/manifestival"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
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

func (h *Helm) Template(ctx context.Context, name, chart string, values runtime.RawExtension) (*release.Release, error) {
	chartRequested, err := loader.LoadDir(chart)
	if err != nil {
		return nil, err
	}

	act := templateAction(name)
	return act.RunWithContext(ctx, chartRequested, map[string]interface{}{})
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
