package event_test

import (
	charles_event "github.com/ZupIT/charlescd-operator/internal/event"
	"github.com/fluxcd/source-controller/api/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

var _ = Describe("RepoStatusPredicate", func() {
	var repoPredicate *charles_event.RepoStatusPredicate
	BeforeEach(func() {
		repoPredicate = charles_event.NewRepoStatusPredicate()
	})
	Context(" when event  is about a  repo resource created", func() {
		It("should return true", func() {
			event := event.CreateEvent{Object: getGitRepository()}
			created := repoPredicate.Create(event)
			Expect(created).To(Equal(true))
		})
	})
	Context(" when event is  not about  a repo resource created", func() {
		It("should return false", func() {
			event := event.CreateEvent{Object: newValidModule()}
			created := repoPredicate.Create(event)
			Expect(created).To(Equal(false))
		})
	})

	Context(" when event is  about  a module resource deleted", func() {
		It("should return false", func() {
			event := event.DeleteEvent{Object: getGitRepository()}
			created := repoPredicate.Delete(event)
			Expect(created).To(Equal(true))
		})
	})

	Context(" when event is  not about  a module resource deleted", func() {
		It("should return false", func() {
			event := event.DeleteEvent{Object: newValidModule()}
			created := repoPredicate.Delete(event)
			Expect(created).To(Equal(false))
		})
	})

	Context(" when event is  about  a repo resource updated", func() {
		It("should return true", func() {
			event := event.UpdateEvent{ObjectNew: getGitRepository(), ObjectOld: new(v1beta1.GitRepository)}
			created := repoPredicate.Update(event)
			Expect(created).To(Equal(true))
		})
	})

	Context(" when event is  not about  a repository resource deleted", func() {
		It("should return false", func() {
			event := event.UpdateEvent{ObjectNew: newValidModule(), ObjectOld: new(v1beta1.GitRepository)}
			created := repoPredicate.Update(event)
			Expect(created).To(Equal(false))
		})
	})

	Context(" when is a generic event about a valid repository", func() {
		It("should return false", func() {
			event := event.GenericEvent{Object: getGitRepository()}
			created := repoPredicate.Generic(event)
			Expect(created).To(Equal(true))
		})
	})

	Context(" when a generic event is  not about a repository resource", func() {
		It("should return false", func() {
			event := event.GenericEvent{Object: newValidModule()}
			created := repoPredicate.Generic(event)
			Expect(created).To(Equal(false))
		})
	})
})

func getGitRepository() *v1beta1.GitRepository {
	return &v1beta1.GitRepository{
		Spec:   v1beta1.GitRepositorySpec{URL: "https://example.com"},
		Status: v1beta1.GitRepositoryStatus{Artifact: &v1beta1.Artifact{URL: "https://example.com/manifests"}},
	}
}
