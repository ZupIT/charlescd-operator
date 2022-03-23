package module_test

import (
	"context"
	charlescdv1alpha1 "github.com/ZupIT/charlescd-operator/api/v1alpha1"
	"github.com/ZupIT/charlescd-operator/pkg/module/mocks"
	"github.com/angelokurtis/reconciler"
	"github.com/fluxcd/pkg/apis/meta"
	"github.com/fluxcd/source-controller/api/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"gopkg.in/h2non/gock.v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/ZupIT/charlescd-operator/pkg/module"
)

func getGitRepositoryWithoutArtifact() *v1beta1.GitRepository {
	condition := metav1.Condition{
		Type:    meta.ReadyCondition,
		Status:  metav1.ConditionFalse,
		Message: "Failed to download artifact",
	}
	return &v1beta1.GitRepository{
		Spec: v1beta1.GitRepositorySpec{URL: manifestLocation},
		Status: v1beta1.GitRepositoryStatus{
			Conditions: []metav1.Condition{condition},
		},
	}
}

var _ = Describe("ArtifactDownload", func() {
	var ctx context.Context
	var statusWriterMock *mocks.StatusWriter
	var gitRepositoryGetter *mocks.GitRepositoryGetter
	var artifactDownload *module.ArtifactDownload

	BeforeEach(func() {
		gock.New("https://example.com").
			Get("/manifests").
			Reply(200).
			BodyString(getArtifactData())
		ctx = context.TODO()
		statusWriterMock = new(mocks.StatusWriter)
		gitRepositoryGetter = new(mocks.GitRepositoryGetter)
		artifactDownload = module.NewArtifactDownload(gitRepositoryGetter, statusWriterMock)
		reconciler.Chain(artifactDownload)
	})
	AfterEach(func() {
		defer gock.Off()
	})

	Context("when reconciling for artifact download", func() {
		It("should update status successfully when source are valid", func() {

			expectedCondition := metav1.Condition{
				Type:    charlescdv1alpha1.SourceReady,
				Status:  metav1.ConditionTrue,
				Reason:  "Downloaded",
				Message: "Sources available locally at",
			}
			mod := setupModule()
			gitRepositoryGetter.On(
				"GetGitRepository",
				mock.Anything, types.NamespacedName{
					Namespace: mod.GetNamespace(),
					Name:      mod.GetName(),
				},
			).Return(getGitRepository(), nil)
			statusWriterMock.On("UpdateModuleStatus", mock.Anything, mod).Return(ctrl.Result{}, nil)

			_, err := artifactDownload.Reconcile(ctx, mod)

			Expect(err).To(BeNil())
			Expect(mod.Status.Conditions[0].Reason).To(Equal(expectedCondition.Reason))
			Expect(mod.Status.Conditions[0].Message).To(ContainSubstring(expectedCondition.Message))
			Expect(mod.Status.Conditions[0].Status).To(Equal(expectedCondition.Status))
			Expect(mod.Status.Conditions[0].Type).To(Equal(expectedCondition.Type))
		})

		It("should update status successfully when artifact is not ready", func() {
			repositoryNotReady := getGitRepositoryWithoutArtifact()
			expectedCondition := metav1.Condition{
				Type:    charlescdv1alpha1.SourceReady,
				Status:  metav1.ConditionFalse,
				Reason:  "Downloaded",
				Message: repositoryNotReady.Status.Conditions[0].Message,
			}
			mod := setupModule()
			gitRepositoryGetter.On(
				"GetGitRepository",
				mock.Anything, types.NamespacedName{
					Namespace: mod.GetNamespace(),
					Name:      mod.GetName(),
				},
			).Return(repositoryNotReady, nil)
			statusWriterMock.On("UpdateModuleStatus", mock.Anything, mod).Return(ctrl.Result{}, nil)

			_, err := artifactDownload.Reconcile(ctx, mod)

			Expect(err).To(BeNil())
			Expect(mod.Status.Conditions[0].Status).To(Equal(expectedCondition.Status))
			Expect(mod.Status.Conditions[0].Type).To(Equal(expectedCondition.Type))
			Expect(mod.Status.Conditions[0].Message).To(Equal(expectedCondition.Message))
		})
	})
})

func getGitRepository() *v1beta1.GitRepository {
	return &v1beta1.GitRepository{
		Spec:   v1beta1.GitRepositorySpec{URL: manifestLocation},
		Status: v1beta1.GitRepositoryStatus{Artifact: &v1beta1.Artifact{URL: "https://example.com/manifests"}},
	}
}
func getArtifactData() string {
	return `
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
}
