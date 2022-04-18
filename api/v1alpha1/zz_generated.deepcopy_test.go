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

package v1alpha1_test

import (
	"github.com/ZupIT/charlescd-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var _ = Describe("ZzGenerated.Deepcopy", func() {
	Context("when copying resources ", func() {
		It("should deep copy a component into another successfully", func() {
			emptyComponent := new(v1alpha1.Component)
			component := newComponent()
			component.DeepCopyInto(emptyComponent)
			Expect(emptyComponent.Name).To(Equal(component.Name))
			Expect(emptyComponent.Containers[0].Name).To(Equal(component.Containers[0].Name))
			Expect(emptyComponent.Containers[0].Image).To(Equal(component.Containers[0].Image))
		})

		It("should copy a component into another successfully", func() {
			component := newComponent()
			copiedComponent := component.DeepCopy()
			Expect(copiedComponent.Name).To(Equal(component.Name))
			Expect(copiedComponent.Containers[0].Name).To(Equal(component.Containers[0].Name))
			Expect(copiedComponent.Containers[0].Image).To(Equal(component.Containers[0].Image))
		})

		It("should copy a container into another successfully", func() {
			container := newComponent().Containers[0]
			containerCopy := new(v1alpha1.Container)
			container.DeepCopyInto(containerCopy)
			Expect(container.Name).To(Equal(containerCopy.Name))
			Expect(container.Image).To(Equal(containerCopy.Image))
		})

		It("should copy a container into another successfully", func() {
			container := newComponent().Containers[0]
			containerCopy := container.DeepCopy()
			Expect(container.Name).To(Equal(containerCopy.Name))
			Expect(container.Image).To(Equal(containerCopy.Image))
		})

		It("should deep copy a gitRef into another successfully", func() {
			gitRef := newGitRef()
			gitRefCopy := new(v1alpha1.GitRef)
			gitRef.DeepCopyInto(gitRefCopy)
			Expect(gitRef.Type).To(Equal(gitRefCopy.Type))
			Expect(gitRef.Value).To(Equal(gitRefCopy.Value))
		})

		It("should copy a git ref into another successfully", func() {
			gitRef := newGitRef()
			gitRefCopy := gitRef.DeepCopy()
			Expect(gitRef.Type).To(Equal(gitRefCopy.Type))
			Expect(gitRef.Value).To(Equal(gitRefCopy.Value))
		})

		It("should deep copy a gitRepository into another successfully", func() {
			gitRepo := newGitRepository()
			gitCopy := new(v1alpha1.GitRepository)
			gitRepo.DeepCopyInto(gitCopy)
			Expect(gitRepo.Path).To(Equal(gitCopy.Path))
			Expect(gitRepo.URL).To(Equal(gitCopy.URL))
			Expect(gitRepo.Interval).To(Equal(gitCopy.Interval))
			Expect(gitRepo.SecretRef.Name).To(Equal(gitCopy.SecretRef.Name))

		})

		It("should copy a git repository into another successfully", func() {
			gitRepo := newGitRepository()
			gitCopy := gitRepo.DeepCopy()
			Expect(gitRepo.Path).To(Equal(gitCopy.Path))
			Expect(gitRepo.URL).To(Equal(gitCopy.URL))
			Expect(gitRepo.Interval).To(Equal(gitCopy.Interval))
			Expect(gitRepo.SecretRef.Name).To(Equal(gitCopy.SecretRef.Name))
		})

		It("should deep copy a helm repository into another successfully", func() {
			helmRepo := newHelmRepository()
			helmCopy := new(v1alpha1.HelmRepository)
			helmRepo.DeepCopyInto(helmCopy)

			Expect(helmRepo.URL).To(Equal(helmCopy.URL))
			Expect(helmRepo.Interval).To(Equal(helmCopy.Interval))
			Expect(helmRepo.SecretRef.Name).To(Equal(helmCopy.SecretRef.Name))
			Expect(helmRepo.HelmChart.Chart).To(Equal(helmCopy.HelmChart.Chart))

		})

		It("should copy a helm repository into another successfully", func() {

			helmRepo := newHelmRepository()
			helmCopy := helmRepo.DeepCopy()

			Expect(helmRepo.URL).To(Equal(helmCopy.URL))
			Expect(helmRepo.Interval).To(Equal(helmCopy.Interval))
			Expect(helmRepo.SecretRef.Name).To(Equal(helmCopy.SecretRef.Name))
			Expect(helmRepo.HelmChart.Chart).To(Equal(helmCopy.HelmChart.Chart))
		})

		It("should deep copy a helm repository into another successfully", func() {
			helm := newHelm()
			helmCopy := new(v1alpha1.Helm)
			helm.DeepCopyInto(helmCopy)

			Expect(helm.GitRepository.Path).To(Equal(helmCopy.GitRepository.Path))
			Expect(helm.HelmRepository.HelmChart.Chart).To(Equal(helmCopy.HelmRepository.HelmChart.Chart))
			Expect(helm.Values.Raw).To(Equal(helmCopy.Values.Raw))

		})

		It("should copy a helm repository into another successfully", func() {

			helm := newHelm()
			helmCopy := helm.DeepCopy()

			Expect(helm.GitRepository.Path).To(Equal(helmCopy.GitRepository.Path))
			Expect(helm.HelmRepository.HelmChart.Chart).To(Equal(helmCopy.HelmRepository.HelmChart.Chart))
			Expect(helm.Values.Raw).To(Equal(helmCopy.Values.Raw))
		})

		It("should deep copy a helm chart into another successfully", func() {
			helmChart := newHelmChart()
			helmCopy := new(v1alpha1.HelmChart)
			helmChart.DeepCopyInto(helmCopy)

			Expect(helmChart.Chart).To(Equal(helmCopy.Chart))

		})

		It("should copy a helm chart into another successfully", func() {

			helmChart := newHelmChart()
			helmCopy := helmChart.DeepCopy()
			Expect(helmChart.Chart).To(Equal(helmCopy.Chart))

		})

		It("should deep copy a kustomization into another successfully", func() {
			kustomization := newKustomization()
			kustomizationCopy := new(v1alpha1.Kustomization)
			kustomization.DeepCopyInto(kustomizationCopy)

			Expect(kustomization.GitRepository.Path).To(Equal(kustomizationCopy.GitRepository.Path))
			Expect(kustomization.Patches.Raw).To(Equal(kustomizationCopy.Patches.Raw))
		})

		It("should copy a kustomization into another successfully", func() {

			kustomization := newKustomization()
			kustomizationCopy := kustomization.DeepCopy()
			Expect(kustomization.GitRepository.Path).To(Equal(kustomizationCopy.GitRepository.Path))
			Expect(kustomization.Patches.Raw).To(Equal(kustomizationCopy.Patches.Raw))
		})

		It("should deep copy a manifest into another successfully", func() {
			manifests := newManifests()
			manifestsCopy := manifests.DeepCopy()

			Expect(manifestsCopy.GitRepository.Path).To(Equal(manifestsCopy.GitRepository.Path))
			Expect(manifestsCopy.GitRepository.Interval).To(Equal(manifestsCopy.GitRepository.Interval))
			Expect(manifestsCopy.GitRepository.URL).To(Equal(manifestsCopy.GitRepository.URL))
		})

		It("should deep copy a module into another successfully", func() {
			module := newModule()
			moduleCopy := module.DeepCopyObject()

			Expect(module.GetObjectKind()).To(Equal(moduleCopy.GetObjectKind()))
		})

		It("should deep copy a module list into another successfully", func() {
			moduleList := newModuleList()
			moduleListCopy := moduleList.DeepCopyObject()

			Expect(moduleList.GetObjectKind()).To(Equal(moduleListCopy.GetObjectKind()))
		})

		It("should deep copy a module list into another successfully", func() {
			moduleList := newModuleList()
			moduleListCopy := moduleList.DeepCopy()

			Expect(moduleList.Items[0].Name).To(Equal(moduleListCopy.Items[0].Name))
		})

		It("should deep copy a module list into another successfully", func() {
			moduleList := newModuleList()
			moduleListCopy := new(v1alpha1.ModuleList)
			moduleList.DeepCopyInto(moduleListCopy)
			Expect(moduleList.Items[0].Name).To(Equal(moduleListCopy.Items[0].Name))
		})

	})
})

func newModuleList() v1alpha1.ModuleList {
	return v1alpha1.ModuleList{Items: []v1alpha1.Module{*newModule()}}
}

func newModule() *v1alpha1.Module {
	module := new(v1alpha1.Module)
	module.Status.Conditions = []metav1.Condition{{Type: "SourceReady", Status: metav1.ConditionTrue}}
	module.Spec.Helm = newHelm()
	return module
}
func newManifests() *v1alpha1.Manifests {
	return &v1alpha1.Manifests{GitRepository: *newGitRepository()}
}

func newKustomization() *v1alpha1.Kustomization {
	return &v1alpha1.Kustomization{
		GitRepository: *newGitRepository(),
		Patches:       &v1.JSON{Raw: []byte("patches")},
	}
}

func newHelmChart() *v1alpha1.HelmChart {
	return &v1alpha1.HelmChart{Chart: "chart"}
}

func newHelm() *v1alpha1.Helm {
	return &v1alpha1.Helm{
		GitRepository:  newGitRepository(),
		HelmRepository: newHelmRepository(),
		Values:         &v1.JSON{Raw: []byte("some-values")},
	}
}

func newHelmRepository() *v1alpha1.HelmRepository {
	return &v1alpha1.HelmRepository{
		Interval:  metav1.Duration{Duration: time.Duration(20)},
		SecretRef: &v1alpha1.SecretRef{Name: "secret"},
		HelmChart: v1alpha1.HelmChart{Chart: "chart"},
		URL:       "https://example.com",
	}
}

func newGitRepository() *v1alpha1.GitRepository {
	return &v1alpha1.GitRepository{Interval: metav1.Duration{Duration: time.Duration(20)},
		Path:      "path",
		URL:       "https://example.com",
		SecretRef: &v1alpha1.SecretRef{Name: "secret-ref"}}
}

func newGitRef() v1alpha1.GitRef {
	return v1alpha1.GitRef{Type: "branch", Value: "main"}
}

func newComponent() v1alpha1.Component {
	container := &v1alpha1.Container{
		Name:  "container",
		Image: "image",
	}
	return v1alpha1.Component{
		Name: "component",
		Containers: []*v1alpha1.Container{
			container,
		},
	}
}
