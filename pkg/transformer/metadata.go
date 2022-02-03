package transformer

import (
	"github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	deployv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
)

type Metadata struct{ reference ObjectReference }

func NewMetadata(reference ObjectReference) *Metadata {
	return &Metadata{reference: reference}
}

func (m *Metadata) TransformMetadata(module *deployv1alpha1.Module) manifestival.Transformer {
	return func(u *unstructured.Unstructured) error {
		u.SetName(module.GetName())
		u.SetNamespace(module.GetNamespace())
		if err := m.reference.SetController(module, u); err != nil {
			return err
		}
		u.SetLabels(map[string]string{
			"app.kubernetes.io/managed-by": "charles",
			"app.kubernetes.io/name":       module.GetName(),
		})
		return nil
	}
}
