package module

import (
	"github.com/fluxcd/pkg/apis/meta"
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type gitRepoStatus struct{ *sourcev1beta1.GitRepository }

func statusOf(gitRepo *sourcev1beta1.GitRepository) *gitRepoStatus {
	return &gitRepoStatus{GitRepository: gitRepo}
}

func (g *gitRepoStatus) IsError() (string, bool) {
	c := apimeta.FindStatusCondition(g.Status.Conditions, meta.ReadyCondition)
	if c == nil {
		return "", false
	}
	if c.Status == metav1.ConditionFalse {
		return c.Message, true
	}
	return "", false
}
