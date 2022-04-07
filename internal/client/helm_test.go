package client_test

import (
	"context"
	"errors"
	client "github.com/ZupIT/charlescd-operator/internal/client"
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

			Expect(err).NotTo(HaveOccurred())
			manifestsReader := &resources.Manifests{}
			helm = client.NewHelm(manifestsReader)
			manifest, err := helm.Template(context.TODO(), "release-name", filepath.Join(currentDir, "./testdata/file.tar.gz"), "fake-app", &v1.JSON{})
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
