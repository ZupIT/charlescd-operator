package eventfilter

import (
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

var repoStatusLog = ctrl.Log.WithName("eventfilter").
	WithValues("name", "gitrepositories", "version", "source.toolkit.fluxcd.io/v1beta1").
	V(1)

type RepoStatus struct{}

func (g *RepoStatus) Create(event event.CreateEvent) bool {
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
	repoStatusLog.Info("Resource created")
	return true
}

func (g *RepoStatus) Update(event event.UpdateEvent) bool {
	if event.ObjectOld == nil || event.ObjectNew == nil {
		return false
	}
	repoOld, ok := event.ObjectOld.(*sourcev1.GitRepository)
	if !ok {
		return false
	}
	repoNew, ok := event.ObjectNew.(*sourcev1.GitRepository)
	if !ok {
		return false
	}
	artifactOld, artifactNew := repoOld.GetArtifact(), repoNew.GetArtifact()
	if artifactOld == nil && artifactNew != nil {
		repoStatusLog.Info("Resource status changed")
		return true
	}
	if artifactOld != nil && artifactNew != nil &&
		artifactOld.Revision != artifactNew.Revision {
		repoStatusLog.Info("Resource status changed")
		return true
	}
	return false
}

func (g *RepoStatus) Delete(event event.DeleteEvent) bool {
	obj := event.Object
	if obj == nil {
		return false
	}
	if _, ok := obj.(*sourcev1.GitRepository); !ok {
		return false
	}
	repoStatusLog.Info("Resource deleted")
	return true
}

func (g *RepoStatus) Generic(event event.GenericEvent) bool {
	return true
}
