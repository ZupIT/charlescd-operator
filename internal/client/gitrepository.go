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

	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ZupIT/charlescd-operator/internal/tracing"
)

type GitRepository struct{ client client.Client }

func NewGitRepository(client client.Client) *GitRepository {
	return &GitRepository{client: client}
}

func (s *GitRepository) GetGitRepository(ctx context.Context, key client.ObjectKey) (*sourcev1beta1.GitRepository, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	m := new(sourcev1beta1.GitRepository)
	err := s.client.Get(ctx, key, m)
	if errors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to lookup resource: %w", err)
	}
	return m, nil
}
