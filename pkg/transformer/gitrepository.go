package transformer

import (
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	charlescdv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
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
