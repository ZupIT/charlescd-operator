package object_test

import (
	"fmt"
	"github.com/ZupIT/charlescd-operator/internal/object"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	//. "github.com/onsi/gomega"
	//"github.com/ZupIT/charlescd-operator/internal/object"
)

var _ = Describe("Converter", func() {
	var converter *object.UnstructuredConverter

	BeforeEach(func() {
		scheme := runtime.NewScheme()
		utilruntime.Must(clientgoscheme.AddToScheme(scheme))
		converter = &object.UnstructuredConverter{Scheme: scheme}
	})

	Context("when converting a resource into another", func() {
		It("should convert a deployment resource into a unstructured", func() {
			var unstructuredDeployment unstructured.Unstructured
			err := converter.ToUnstructured(getDeploymentResource(), &unstructuredDeployment)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should convert a unstructured resource into a deployment", func() {
			var deployment appsv1.Deployment
			err := converter.FromUnstructured(getUnstructuredDeployment(), &deployment)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error when fails to convert a resource from a unstructured", func() {
			msgError := `error converting from unstructured object: no kind "Deployment" is registered for version ` +
				`"apps/v1" in scheme "pkg/runtime/scheme.go:100"`
			var deployment appsv1.Deployment
			converter.Scheme = runtime.NewScheme()
			err := converter.FromUnstructured(getUnstructuredDeployment(), &deployment)
			Expect(err.Error()).To(Equal(msgError))
		})

		It("should convert a deployment resource into a unstructured", func() {
			msgError := `error converting to unstructured object: no kind is registered for the type v1.Deployment ` +
				`in scheme "pkg/runtime/scheme.go:100"`
			var unstructuredDeployment unstructured.Unstructured
			converter.Scheme = runtime.NewScheme()
			err := converter.ToUnstructured(getDeploymentResource(), &unstructuredDeployment)
			fmt.Println(err.Error())
			Expect(err.Error()).To(Equal(msgError))
		})
	})

})

func getDeploymentResource() *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

func getUnstructuredDeployment() *unstructured.Unstructured {
	deployment := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":      "test-deploy",
				"namespace": "test",
			},
		},
	}
	return deployment
}
