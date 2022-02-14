package event

import (
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type RepoStatusPredicate struct{}

func NewRepoStatusPredicate() *RepoStatusPredicate {
	return &RepoStatusPredicate{}
}

func (g *RepoStatusPredicate) Create(event event.CreateEvent) bool {
	obj := event.Object
	if obj == nil {
		return false
	}
	repo, ok := obj.(*sourcev1.GitRepository)
	if !ok {
		return false
	}
	artifact := repo.GetArtifact()
	if artifact == nil {
		return false
	}
	log.WithValues("name", obj.GetName(),
		"namespace", obj.GetNamespace(),
		"resourceVersion", obj.GetResourceVersion()).
		Info("a new GitRepository has been created")
	return true
}

func (g *RepoStatusPredicate) Update(event event.UpdateEvent) bool {
	objOld, objNew := event.ObjectOld, event.ObjectNew
	if objOld == nil || objNew == nil {
		return false
	}
	repoOld, ok := objOld.(*sourcev1.GitRepository)
	if !ok {
		return false
	}
	repoNew, ok := objNew.(*sourcev1.GitRepository)
	if !ok {
		return false
	}
	artifactOld, artifactNew := repoOld.GetArtifact(), repoNew.GetArtifact()
	if artifactOld == nil && artifactNew != nil {
		log.WithValues("name", objNew.GetName(),
			"namespace", objNew.GetNamespace(),
			"resourceVersion", objNew.GetResourceVersion(),
			"diff", diff(
				&sourcev1.GitRepository{Spec: repoOld.Spec, Status: repoOld.Status},
				&sourcev1.GitRepository{Spec: repoNew.Spec, Status: repoNew.Status})).
			Info("a GitRepository was updated")
		return true
	}
	if artifactOld != nil && artifactNew != nil &&
		artifactOld.Revision != artifactNew.Revision {
		log.WithValues("name", objNew.GetName(),
			"namespace", objNew.GetNamespace(),
			"resourceVersion", objNew.GetResourceVersion(),
			"diff", diff(
				&sourcev1.GitRepository{Spec: repoOld.Spec, Status: repoOld.Status},
				&sourcev1.GitRepository{Spec: repoNew.Spec, Status: repoNew.Status})).
			Info("a GitRepository was updated")
		return true
	}
	return false
}

func (g *RepoStatusPredicate) Delete(event event.DeleteEvent) bool {
	obj := event.Object
	if obj == nil {
		return false
	}
	if _, ok := obj.(*sourcev1.GitRepository); !ok {
		return false
	}
	log.WithValues("name", obj.GetName(),
		"namespace", obj.GetNamespace(),
		"resourceVersion", obj.GetResourceVersion()).
		Info("a GitRepository was deleted")
	return true
}

func (g *RepoStatusPredicate) Generic(event event.GenericEvent) bool {
	return true
}
