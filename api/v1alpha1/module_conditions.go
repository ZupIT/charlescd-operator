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
		Message: "Sources available locally at " + path,
	})
	in.Status.Source = &Source{Path: path}
	in.updatePhase()
	return updated(old, in)
}

func (in *Module) SetSourceValid() bool {
	old := in.DeepCopy()
	meta.SetStatusCondition(&in.Status.Conditions, metav1.Condition{
		Type:    SourceValid,
		Status:  metav1.ConditionTrue,
		Reason:  "Validated",
		Message: "Helm chart templates were successfully rendered",
	})
	in.updatePhase()
	return updated(old, in)
}

func (in *Module) SetSourceError(reason, message string) bool {
	old := in.DeepCopy()
	meta.SetStatusCondition(&in.Status.Conditions, metav1.Condition{
		Type:    SourceReady,
		Status:  metav1.ConditionFalse,
		Reason:  reason,
		Message: message,
	})
	meta.RemoveStatusCondition(&in.Status.Conditions, SourceValid)
	in.Status.Source = nil
	in.Status.Components = nil
	in.updatePhase()
	return updated(old, in)
}

func (in *Module) SetSourceInvalid(reason, message string) bool {
	old := in.DeepCopy()
	meta.SetStatusCondition(&in.Status.Conditions, metav1.Condition{
		Type:    SourceValid,
		Status:  metav1.ConditionFalse,
		Reason:  reason,
		Message: message,
	})
	in.Status.Components = nil
	in.updatePhase()
	return updated(old, in)
}

func (in *Module) SetComponents(components []*Component) bool {
	old := in.DeepCopy()
	in.Status.Components = components
	return updated(old, in)
}

func (in *Module) UpdatePhase() bool {
	old := in.DeepCopy()
	in.updatePhase()
	return updated(old, in)
}

func (in *Module) updatePhase() {
	switch conditions := in.Status.Conditions; {
	case meta.IsStatusConditionTrue(conditions, SourceReady) && meta.IsStatusConditionTrue(conditions, SourceValid):
		in.Status.Phase = "Ready"
	case meta.IsStatusConditionFalse(conditions, SourceReady) || meta.IsStatusConditionFalse(conditions, SourceValid):
		in.Status.Phase = "Failed"
	default:
		in.Status.Phase = "Processing"
	}
}

func updated(oldest, newest client.Object) bool {
	patch := client.MergeFrom(oldest)
	data, _ := patch.Data(newest)
	return string(data) != "{}"
}
