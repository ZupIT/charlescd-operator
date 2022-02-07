package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (in *Module) SetSourceReady(path string) (string, bool) {
	old := in.DeepCopy()
	meta.SetStatusCondition(&in.Status.Conditions, metav1.Condition{
		Type:    SourceReady,
		Status:  metav1.ConditionTrue,
		Reason:  "Downloaded",
		Message: "Sources are available at: " + path,
	})
	in.Status.Source = &Source{Path: path}
	in.updatePhase()
	return updated(old, in)
}

func (in *Module) SetSourceError(reason, message string) (string, bool) {
	old := in.DeepCopy()
	meta.SetStatusCondition(&in.Status.Conditions, metav1.Condition{
		Type:    SourceReady,
		Status:  metav1.ConditionFalse,
		Reason:  reason,
		Message: message,
	})
	in.Status.Source = nil
	in.updatePhase()
	return updated(old, in)
}

func (in *Module) RemoveSource() (string, bool) {
	old := in.DeepCopy()
	meta.RemoveStatusCondition(&in.Status.Conditions, SourceReady)
	in.Status.Source = nil
	in.updatePhase()
	return updated(old, in)
}

func (in *Module) UpdatePhase() (string, bool) {
	old := in.DeepCopy()
	in.updatePhase()
	return updated(old, in)
}

func (in *Module) updatePhase() {
	switch conditions := in.Status.Conditions; {
	case meta.IsStatusConditionTrue(conditions, SourceReady):
		in.Status.Phase = "Ready"
	case meta.IsStatusConditionFalse(conditions, SourceReady):
		in.Status.Phase = "Failed"
	default:
		in.Status.Phase = "Processing"
	}
}

func updated(old, new client.Object) (diff string, updated bool) {
	patch := client.MergeFrom(old)
	data, _ := patch.Data(new)
	diff = string(data)
	updated = diff != "{}"
	return diff, updated
}
