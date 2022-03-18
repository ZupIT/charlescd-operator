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

package module_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/ZupIT/charlescd-operator/pkg/module"
	"github.com/ZupIT/charlescd-operator/test/module/mocks"
	"github.com/angelokurtis/reconciler"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"
)

var manifestLocation = filepath.Join(os.TempDir(), "deployment.yaml")

var _ = Describe("Module", func() {
	var ctx context.Context
	var statusWriterMock *mocks.StatusWriter
	var manifestClientMock *mocks.ManifestClient
	var manifestValidation *module.ManifestValidation

	BeforeEach(func() {
		ctx = context.TODO()
		statusWriterMock = new(mocks.StatusWriter)
		manifestClientMock = new(mocks.ManifestClient)
		manifestValidation = module.NewManifestValidation(statusWriterMock, manifestClientMock)
		reconciler.Chain(manifestValidation)
	})
	Context("when reconciling  pure manifests", func() {
		It("should update status successfully when source are valid", func() {
			expectedCondition := metav1.Condition{
				Type:    charlescdv1alpha1.SourceValid,
				Status:  metav1.ConditionTrue,
				Reason:  "Validated",
				Message: module.SuccessManifestLoadMessage,
			}
			storeManifestInFSys()
			mod := setupModule()

			manifestClientMock.On(
				"DownloadFromSource",
				mock.Anything, mod.Status.Source.Path, mod.Spec.Manifests.GitRepository.Path,
			).Return(manifestLocation, nil)
			statusWriterMock.On("UpdateModuleStatus", mock.Anything, mod).Return(ctrl.Result{}, nil)

			_, err := manifestValidation.Reconcile(ctx, mod)

			Expect(err).To(BeNil())
			Expect(mod.Status.Conditions[1].Reason).To(Equal(expectedCondition.Reason))
			Expect(mod.Status.Conditions[1].Message).To(Equal(expectedCondition.Message))
			Expect(mod.Status.Conditions[1].Status).To(Equal(expectedCondition.Status))
			Expect(mod.Status.Conditions[1].Type).To(Equal(expectedCondition.Type))
		})

		It("should update status correctly when fails to download resource", func() {
			downloadError := errors.New("failed to download from source")
			mod := setupModule()

			manifestClientMock.On("DownloadFromSource", mock.Anything,
				mod.Status.Source.Path,
				mod.Spec.Manifests.GitRepository.Path,
			).Return("", downloadError)
			statusWriterMock.On("UpdateModuleStatus", mock.Anything, mod).Return(ctrl.Result{}, nil)

			_, err := manifestValidation.Reconcile(ctx, mod)

			Expect(err).To(BeNil())
			Expect(mod.Status.Conditions[1].Type).To(Equal(charlescdv1alpha1.SourceValid))
			Expect(mod.Status.Conditions[1].Status).To(Equal(metav1.ConditionFalse))
			Expect(mod.Status.Conditions[1].Message).To(Equal(downloadError.Error()))
		})

		It("should return error when fails to load manifests", func() {
			mod := setupModule()

			manifestClientMock.On(
				"DownloadFromSource",
				mock.Anything,
				mod.Status.Source.Path,
				mod.Spec.Manifests.GitRepository.Path,
			).Return("./fake-path/deployment.yaml", nil)
			statusWriterMock.On("UpdateModuleStatus", mock.Anything, mod).Return(ctrl.Result{}, nil)

			_, err := manifestValidation.Reconcile(ctx, mod)

			Expect(err).To(BeNil())
			Expect(mod.Status.Conditions[1].Type).To(Equal(charlescdv1alpha1.SourceValid))
			Expect(mod.Status.Conditions[1].Status).To(Equal(metav1.ConditionFalse))
			Expect(mod.Status.Conditions[1].Message).To(Equal("stat ./fake-path/deployment.yaml: no such file or directory"))
		})
	})
})

func setupModule() *charlescdv1alpha1.Module {
	module := new(charlescdv1alpha1.Module)
	module.Status.Conditions = []metav1.Condition{{Type: "SourceReady", Status: metav1.ConditionTrue}}
	module.Spec.Manifests = &charlescdv1alpha1.Manifests{GitRepository: charlescdv1alpha1.GitRepository{URL: "https://example.com"}}
	module.Status.Source = &charlescdv1alpha1.Source{Path: "path/file.tgz"}
	return module
}

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
