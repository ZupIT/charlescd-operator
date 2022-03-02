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
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter"
	mf "github.com/manifestival/manifestival"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/kustomize/v4/commands/build"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

type Kustomization struct {
	manifests ManifestsReader
}

func NewKustomization(manifests ManifestsReader) *Kustomization {
	return &Kustomization{manifests: manifests}
}

func (k *Kustomization) Kustomize(ctx context.Context, source, path string) (mf.Manifest, error) {
	destination := source[0 : len(source)-len(".tar.gz")]
	if path != "" {
		source += "//" + path
		destination = filepath.Join(destination, path)
	}

	if err := getter.GetAny(destination, source); err != nil {
		return mf.Manifest{}, fmt.Errorf("error extracting Source artifact %s: %w", source, err)
	}
	defer os.RemoveAll(destination)
	resMap, err := k.kustomize(destination)
	if err != nil {
		return mf.Manifest{}, err
	}
	r, err := resMap.AsYaml()
	if err != nil {
		return mf.Manifest{}, fmt.Errorf("erro on resMap marshal: %w", err)
	}

	return k.manifests.FromBytes(ctx, r)
}

func (k *Kustomization) kustomize(filePath string) (resmap.ResMap, error) {
	kustomizer := krusty.MakeKustomizer(
		build.HonorKustomizeFlags(krusty.MakeDefaultOptions()))

	fsys := filesys.MakeFsOnDisk()

	resMap, err := kustomizer.Run(fsys, filePath)
	if err != nil {
		return resMap, fmt.Errorf("erro running kustomizer: %w", err)
	}

	return resMap, nil
}
