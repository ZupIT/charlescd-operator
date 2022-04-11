package event_test

import (
	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	charles_event "github.com/ZupIT/charlescd-operator/internal/event"
	"github.com/fluxcd/source-controller/api/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

var _ = Describe("ModulePredicate", func() {
	var modulePredicate *charles_event.ModulePredicate
	BeforeEach(func() {
		modulePredicate = charles_event.NewModulePredicate()
	})
	Context(" when event  is about a  module resource created", func() {
		It("should return true", func() {
			event := event.CreateEvent{Object: newValidModule()}
			created := modulePredicate.Create(event)
			Expect(created).To(Equal(true))
		})
	})
	Context(" when event is  not about  a module resource created", func() {
		It("should return false", func() {
			event := event.CreateEvent{Object: new(v1beta1.GitRepository)}
			created := modulePredicate.Create(event)
			Expect(created).To(Equal(false))
		})
	})

	Context(" when event is  about  a module resource deleted", func() {
		It("should return false", func() {
			event := event.DeleteEvent{Object: newValidModule()}
			created := modulePredicate.Delete(event)
			Expect(created).To(Equal(true))
		})
	})

	Context(" when event is  not about  a module resource deleted", func() {
		It("should return false", func() {
			event := event.DeleteEvent{Object: new(v1beta1.GitRepository)}
			created := modulePredicate.Delete(event)
			Expect(created).To(Equal(false))
		})
	})

	Context(" when event is  about  a module resource ", func() {
		It("should return false", func() {
			event := event.UpdateEvent{ObjectNew: newValidModule(), ObjectOld: newValidModule()}
			created := modulePredicate.Update(event)
			Expect(created).To(Equal(true))
		})
	})

	Context(" when event is  not about  a module resource deleted", func() {
		It("should return false", func() {
			event := event.UpdateEvent{ObjectNew: new(v1beta1.GitRepository), ObjectOld: new(v1beta1.GitRepository)}
			created := modulePredicate.Update(event)
			Expect(created).To(Equal(false))
		})
	})

	Context(" when is a generic event about a valid module", func() {
		It("should return false", func() {
			event := event.GenericEvent{Object: newValidModule()}
			created := modulePredicate.Generic(event)
			Expect(created).To(Equal(true))
		})
	})

	Context(" when a generic event is  not about a module resource", func() {
		It("should return false", func() {
			event := event.GenericEvent{Object: new(v1beta1.GitRepository)}
			created := modulePredicate.Generic(event)
			Expect(created).To(Equal(false))
		})
	})

})

func newValidModule() *charlescdv1alpha1.Module {
	module := new(charlescdv1alpha1.Module)
	module.Status.Conditions = []metav1.Condition{{Type: "SourceReady", Status: metav1.ConditionTrue}, {Type: "SourceValid", Status: metav1.ConditionTrue}}
	module.Spec.Manifests = &charlescdv1alpha1.Manifests{GitRepository: charlescdv1alpha1.GitRepository{URL: "https://example.com"}}
	module.Status.Source = &charlescdv1alpha1.Source{Path: "path/file.tgz"}
	module.SetNamespace("test")
	module.SetName("test-module")
	return module
}
