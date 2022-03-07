package event

import (
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type RepoStatusPredicate struct{}

func NewRepoStatusPredicate() *RepoStatusPredicate {
	return &RepoStatusPredicate{}
}

func (g *RepoStatusPredicate) Create(event.CreateEvent) bool { return false }

func (g *RepoStatusPredicate) Update(event event.UpdateEvent) bool {
	objOld, objNew := event.ObjectOld, event.ObjectNew
	if objOld == nil || objNew == nil {
		return false
	}
	repoOld, ok := objOld.(*sourcev1beta1.GitRepository)
	if !ok {
		return false
	}
	repoNew, ok := objNew.(*sourcev1beta1.GitRepository)
	if !ok {
		return false
	}
	artifactOld, artifactNew := repoOld.GetArtifact(), repoNew.GetArtifact()
	if artifactOld == nil && artifactNew != nil {
		log.WithValues("name", objNew.GetName(),
			"namespace", objNew.GetNamespace(),
			"resourceVersion", objNew.GetResourceVersion(),
			"diff", diff(
				&sourcev1beta1.GitRepository{Spec: repoOld.Spec, Status: repoOld.Status},
				&sourcev1beta1.GitRepository{Spec: repoNew.Spec, Status: repoNew.Status})).
			Info("GitRepository updated")
		return true
	}
	if artifactOld != nil && artifactNew != nil &&
		artifactOld.Revision != artifactNew.Revision {
		log.WithValues("name", objNew.GetName(),
			"namespace", objNew.GetNamespace(),
			"resourceVersion", objNew.GetResourceVersion(),
			"diff", diff(
				&sourcev1beta1.GitRepository{Spec: repoOld.Spec, Status: repoOld.Status},
				&sourcev1beta1.GitRepository{Spec: repoNew.Spec, Status: repoNew.Status})).
			Info("GitRepository updated")
		return true
	}
	return false
}

func (g *RepoStatusPredicate) Delete(event event.DeleteEvent) bool {
	obj := event.Object
	if obj == nil {
		return false
	}
	if _, ok := obj.(*sourcev1beta1.GitRepository); !ok {
		return false
	}
	log.WithValues("name", obj.GetName(),
		"namespace", obj.GetNamespace(),
		"resourceVersion", obj.GetResourceVersion()).
		Info("GitRepository deleted")
	return true
}

func (g *RepoStatusPredicate) Generic(event.GenericEvent) bool { return false }
