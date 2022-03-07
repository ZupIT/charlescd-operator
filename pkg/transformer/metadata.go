package transformer

import (
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	charlescdv1alpha1 "github.com/tiagoangelozup/charles-alpha/api/v1alpha1"
)

type Metadata struct{ reference ObjectReference }

func NewMetadata(reference ObjectReference) *Metadata {
	return &Metadata{reference: reference}
}

func (m *Metadata) TransformMetadata(module *charlescdv1alpha1.Module) mf.Transformer {
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
