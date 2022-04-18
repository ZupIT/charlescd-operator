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

package event

import (
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
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
	if _, ok := obj.(*sourcev1beta1.GitRepository); !ok {
		return false
	}
	log.WithValues("name", obj.GetName(),
		"namespace", obj.GetNamespace(),
		"resourceVersion", obj.GetResourceVersion()).
		Info("GitRepository created")
	return true
}

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

func (g *RepoStatusPredicate) Generic(event event.GenericEvent) bool {
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
		Info("GitRepository generic")
	return true
}
