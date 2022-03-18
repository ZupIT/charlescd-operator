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
)

type Manifest struct{}

func NewManifest(manifests ManifestsReader) *Manifest {
	return &Manifest{}
}

func (h *Manifest) DownloadFromSource(ctx context.Context, source, path string) (string, error) {
	destination, err := os.MkdirTemp(os.TempDir(), "manifests")
	if err != nil {
		return "", fmt.Errorf("error creating temp dir %w", err)
	}
	if err := getter.GetAny(destination, source); err != nil {
		return "", fmt.Errorf("error downloading manifests  %w", err)
	}
	if path != "" {
		return filepath.Join(destination, path), nil
	}
	return destination, nil
}
