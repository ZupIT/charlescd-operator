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

var manifestLocation = filepath.Join(os.TempDir(), "deployment.yaml")

var _ = Describe("Manifest", Ordered, func() {

	BeforeAll(func() {
		storeManifestInFSys()
	})

	Context("When loading manifests", func() {
		It("should load correctly from source", func() {
			manifests := &resources.Manifests{}

			manifestClient := client.NewManifest(manifests)

			mf, err := manifestClient.LoadFromSource(context.TODO(), manifestLocation, "")

			Expect(err).To(BeNil())
			loadedManifests := mf.Resources()
			Expect(len(loadedManifests)).To(Equal(1))
			Expect(loadedManifests[0].GetKind()).To(Equal("Deployment"))
			Expect(loadedManifests[0].GetAPIVersion()).To(Equal("apps/v1"))
		})
	})

	AfterAll(func() {
		err := os.RemoveAll(manifestLocation)
		Expect(err).To(BeNil())
	})

})

func storeManifestInFSys() {
	fileData := `
# Source: event-receiver/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: event-receiver
  labels:
    helm.sh/chart: event-receiver-0.1.0
    app.kubernetes.io/name: event-receiver
    app.kubernetes.io/instance: RELEASE-NAME
    app.kubernetes.io/version: "1.16.0"
    app.kubernetes.io/managed-by: Helm
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: event-receiver
      app.kubernetes.io/instance: RELEASE-NAME
  template:
    metadata:
      labels:
        app.kubernetes.io/name: event-receiver
        app.kubernetes.io/instance: RELEASE-NAME
    spec:
      serviceAccountName: event-receiver
      securityContext:
        {}
      containers:
        - name: event-receiver
          securityContext:
            {}
          image: "thallesf/event-receiver:1.0"
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 3000
              protocol: TCP
          livenessProbe:
            httpGet:
              path: /
              port: http
          readinessProbe:
            httpGet:
              path: /
              port: http
          resources:
            {}`
	f, err := os.Create(manifestLocation)
	Expect(err).To(BeNil())
	_, err = f.WriteString(fileData)
	Expect(err).To(BeNil())
	defer f.Close()
}
