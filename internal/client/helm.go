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

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter"
	mf "github.com/manifestival/manifestival"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
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

func (h *Helm) Template(ctx context.Context, releaseName, source, path string, values *apiextensionsv1.JSON) (mf.Manifest, error) {
	destination := source[0 : len(source)-len(".tar.gz")]
	if path != "" {
		source += "//" + path
		destination = filepath.Join(destination, path)
	}

	if err := getter.GetAny(destination, source); err != nil {
		return mf.Manifest{}, fmt.Errorf("error extracting Source artifact %s: %w", source, err)
	}
	defer os.RemoveAll(destination)

	var v map[string]interface{}
	if values != nil {
		_ = json.Unmarshal(values.Raw, &v)
	}

	r, err := h.template(ctx, releaseName, destination, v)
	if err != nil {
		return mf.Manifest{}, err
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

	r, err := client.RunWithContext(ctx, c, values)
	if err != nil {
		return nil, fmt.Errorf("error rendering Helm chart templates: %w", err)
	}
	return r, nil
}
