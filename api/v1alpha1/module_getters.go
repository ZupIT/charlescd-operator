package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
)

func (in *Module) GetGitRepository() (*GitRepository, error) {
	repos := make([]*GitRepository, 0, 0)
	if in.Spec.Helm != nil && in.Spec.Helm.GitRepository != nil {
		repos = append(repos, in.Spec.Helm.GitRepository)
	}
	if in.Spec.Kustomization != nil {
		repos = append(repos, &in.Spec.Kustomization.GitRepository)
	}
	if in.Spec.Manifests != nil {
		repos = append(repos, &in.Spec.Manifests.GitRepository)
	}
	total := len(repos)
	if total > 1 {
		return nil, fmt.Errorf("invalid module definition: expected 1 GitRepository, got %d", total)
	}
	if total == 1 {
		return repos[0], nil
	}
	return nil, nil
}

func (in *Module) IsSourceReady() bool {
	return meta.IsStatusConditionTrue(in.Status.Conditions, SourceReady)
}

func (in *Module) IsSourceValid() bool {
	return meta.IsStatusConditionTrue(in.Status.Conditions, SourceValid)
}
