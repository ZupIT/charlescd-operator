/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ModuleSpec defines the desired state of Module
type ModuleSpec struct {
	Repository Repository           `json:"repository"`
	Values     runtime.RawExtension `json:"values,omitempty"`
}

// Repository defines the location where sources are stored
type Repository struct {
	Git  *Git  `json:"git,omitempty"`
	Helm *Helm `json:"helm,omitempty"`
}

// Git defines the address where sources are tracked
type Git struct {
	URL string `json:"url,omitempty"`
	Ref *Ref   `json:"ref,omitempty"`
}

// Ref defines references to a specific commit
type Ref struct {
	Branch string `json:"branch,omitempty"`
	Commit string `json:"commit,omitempty"`
	Tag    string `json:"tag,omitempty"`
}

// Helm defines the address where charts are packaged
type Helm struct {
	URL     string `json:"url,omitempty"`
	Chart   string `json:"chart,omitempty"`
	Version string `json:"version,omitempty"`
}

// ModuleStatus defines the observed state of Module
type ModuleStatus struct { // INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Module is the Schema for the modules API
type Module struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModuleSpec   `json:"spec,omitempty"`
	Status ModuleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ModuleList contains a list of Module
type ModuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Module `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Module{}, &ModuleList{})
}
