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

	mf "github.com/manifestival/manifestival"
)

type ManifestsReader interface {
	FromString(ctx context.Context, manifests string) (mf.Manifest, error)
	FromBytes(ctx context.Context, manifests []byte) (mf.Manifest, error)
	FromPath(ctx context.Context, path string, recursive bool) (mf.Manifest, error)
}
