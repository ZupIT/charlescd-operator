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

package transformer

import (
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
)

type GitRepository struct{ object ObjectConverter }

func NewGitRepository(object ObjectConverter) *GitRepository {
	return &GitRepository{object: object}
}

func (g *GitRepository) TransformGitRepository(module *charlescdv1alpha1.Module) mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		git, err := module.GetGitRepository()
		if err != nil {
			return err
		}
		if u.GetKind() == "GitRepository" && git != nil {
			gitrepo := &sourcev1beta1.GitRepository{}
			if err = g.object.FromUnstructured(u, gitrepo); err != nil {
				return err
			}
			gitrepo.Spec.URL = git.URL
			switch git.Ref.Type {
			case "branch":
				gitrepo.Spec.Reference = &sourcev1beta1.GitRepositoryRef{Branch: git.Ref.Value}
			case "tag":
				gitrepo.Spec.Reference = &sourcev1beta1.GitRepositoryRef{Tag: git.Ref.Value}
			case "commit":
				gitrepo.Spec.Reference = &sourcev1beta1.GitRepositoryRef{Commit: git.Ref.Value}
			case "semver":
				gitrepo.Spec.Reference = &sourcev1beta1.GitRepositoryRef{SemVer: git.Ref.Value}
			}
			if err = g.object.ToUnstructured(gitrepo, u); err != nil {
				return err
			}
		}
		return nil
	}
}
