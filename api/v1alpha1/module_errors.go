package v1alpha1

import "fmt"

//+kubebuilder:object:generate=false
type MultipleGitRepositoryError struct{ expected, got int }

func (e *MultipleGitRepositoryError) Error() string {
	return fmt.Sprintf("invalid module definition: expected %d GitRepository, got %d", e.expected, e.got)
}
