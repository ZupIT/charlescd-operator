package resources_test

import (
	"context"
	"github.com/ZupIT/charlescd-operator/internal/resources"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manifests.", func() {
	var manifests resources.Manifests
	var ctx context.Context
	BeforeEach(func() {
		manifests = resources.Manifests{}
		ctx = context.TODO()
	})

	Context("when loading manifests", func() {
		It("should load successfully from defaults ", func() {
			mf, err := manifests.LoadDefaults(ctx)
			Expect(err).ToNot(HaveOccurred())
			for _, resource := range mf.Resources() {
				Expect(resource.GetKind()).To(Equal("GitRepository"))
				Expect(resource.GetName()).To(Equal("default"))
			}
		})

		It("should load successfully from bytes", func() {
			mf, err := manifests.FromBytes(context.TODO(), getArtifactDataBytes())
			Expect(err).ToNot(HaveOccurred())
			for _, resource := range mf.Resources() {
				Expect(resource.GetName()).To(Equal("quiz-app"))
				Expect(resource.GetKind()).To(Equal("Deployment"))
			}
		})

		It("should load successfully from bytes", func() {
			mf, err := manifests.FromString(context.TODO(), getArtifactData())
			Expect(err).ToNot(HaveOccurred())
			for _, resource := range mf.Resources() {
				Expect(resource.GetName()).To(Equal("quiz-app"))
				Expect(resource.GetKind()).To(Equal("Deployment"))
			}
		})
	})
})

func getArtifactDataBytes() []byte {
	return []byte(getArtifactData())
}
func getArtifactData() string {
	return `
# Source: event-receiver/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: quiz-app
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
        - name: quiz-app
          securityContext:
            {}
          image: "charlescd/quiz-app:1.0"
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
}
