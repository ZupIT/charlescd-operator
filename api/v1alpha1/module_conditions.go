package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (in *Module) SetSourceReady(path string) bool {
	old := in.DeepCopy()
	meta.SetStatusCondition(&in.Status.Conditions, metav1.Condition{
		Type:    SourceReady,
		Status:  metav1.ConditionTrue,
		Reason:  "Downloaded",
		Message: "Sources are available at: " + path,
	})
	in.Status.Source = &Source{Path: path}
	return updated(old, in)
}

func (in *Module) SetSourceError(reason string, err error) bool {
	old := in.DeepCopy()
	meta.SetStatusCondition(&in.Status.Conditions, metav1.Condition{
		Type:    SourceReady,
		Status:  metav1.ConditionFalse,
		Reason:  reason,
		Message: err.Error(),
	})
	in.Status.Source = nil
	return updated(old, in)
}

func updated(old, new client.Object) bool {
	patch := client.MergeFrom(old)
	data, _ := patch.Data(new)
	return string(data) != "{}"
}
