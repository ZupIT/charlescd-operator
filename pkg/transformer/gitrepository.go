package transformer

import (
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
)

type GitRepository struct{ object ObjectConverter }

func NewGitRepository(object ObjectConverter) *GitRepository {
	return &GitRepository{object: object}
}

func (g *GitRepository) TransformGitRepository(module *deployv1alpha1.Module) mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		git := module.Spec.Repository.Git
		if u.GetKind() == "GitRepository" && git != nil {
			gitrepo := &sourcev1.GitRepository{}
			if err := g.object.FromUnstructured(u, gitrepo); err != nil {
				return err
			}
			gitrepo.Spec.URL = git.URL
			gitrepo.Spec.Reference = &sourcev1.GitRepositoryRef{
				Branch: git.Branch,
				Tag:    git.Tag,
				Commit: git.Commit,
			}
			if err := g.object.ToUnstructured(gitrepo, u); err != nil {
				return err
			}
		}
		return nil
	}
}
