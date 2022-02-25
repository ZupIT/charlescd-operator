package filter

import (
	mf "github.com/manifestival/manifestival"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
)

type GitRepository struct{}

func (g *GitRepository) FilterGitRepository(module *deployv1alpha1.Module) mf.Predicate {
	if git, _ := module.GetGitRepository(); git == nil {
		return mf.Not(mf.ByKind("GitRepository"))
	}
	return mf.Everything
}
