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

package client_test

import (
	"context"
	"github.com/ZupIT/charlescd-operator/internal/client"
	"github.com/ZupIT/charlescd-operator/internal/resources"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("Kustomization", func() {
	Context("when are valid kustomize manifests", func() {
		It("should render successfully", func() {
			currentDir, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())
			source := filepath.Join(currentDir, "./testdata/kustomize.tar.gz")
			manifestsReader := &resources.Manifests{}
			kustomization := client.NewKustomization(manifestsReader)
			manifest, err := kustomization.Kustomize(context.TODO(), source, "./internal/client/testdata/kustomize")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(manifest.Resources())).To(Equal(2))
		})
	})

	Context("when the path for kustomize manifests does not exists", func() {
		It("should return error ", func() {
			currentDir, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())
			source := filepath.Join(currentDir, "./testdata/kustomization.tar.gz")
			manifestsReader := &resources.Manifests{}
			kustomization := client.NewKustomization(manifestsReader)
			_, err = kustomization.Kustomize(context.TODO(), source, "./internal/client/testdata/kustomize")
			Expect(err).To(HaveOccurred())

		})
	})

})
