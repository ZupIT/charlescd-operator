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
	"errors"
	"github.com/ZupIT/charlescd-operator/internal/client"
	"github.com/ZupIT/charlescd-operator/internal/resources"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"os"
	"path/filepath"
)

var _ = Describe("Helm", func() {

	var helm *client.Helm

	Context("when chart is valid", func() {

		It("should render manifests successfully", func() {
			currentDir, err := os.Getwd()
			source := filepath.Join(currentDir, "./testdata/helm/chart.tar.gz")
			manifestsReader := &resources.Manifests{}
			helm = client.NewHelm(manifestsReader)
			manifest, err := helm.Template(context.TODO(), "release-name", source, "./internal/client/testdata/helm/fake-app", &v1.JSON{})
			Expect(err).ToNot(HaveOccurred())
			resources := manifest.Resources()
			Expect(len(resources)).To(Equal(2))
			Expect(resources[0].GetKind()).To(Equal("Service"))
			Expect(resources[1].GetKind()).To(Equal("Deployment"))
		})

		It("should return error", func() {
			expectedError := errors.New("error extracting Source artifact /fake-path: stat /fake-path: no such file or directory")
			manifestsReader := &resources.Manifests{}
			helm = client.NewHelm(manifestsReader)
			_, err := helm.Template(context.TODO(), "release-name", "/fake-path", "", &v1.JSON{})
			Expect(err.Error()).To(Equal(expectedError.Error()))
		})
	})

})
