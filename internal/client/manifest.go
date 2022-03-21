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
)

type Manifest struct {
	reader ManifestsReader
}

func NewManifest(reader ManifestsReader) *Manifest {
	return &Manifest{reader}
}

func (h *Manifest) LoadFromSource(ctx context.Context, source, path string) (mf.Manifest, error) {
	destination, err := os.MkdirTemp(os.TempDir(), "manifests")
	if err != nil {
		return mf.Manifest{}, fmt.Errorf("error creating temp dir %w", err)
	}
	if err := getter.GetAny(destination, source); err != nil {
		return mf.Manifest{}, fmt.Errorf("error downloading manifests  %w", err)
	}
	defer os.RemoveAll(destination)
	if path != "" {
		return h.reader.FromPath(ctx, filepath.Join(destination, path))
	}
	return h.reader.FromPath(ctx, destination)
}
