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
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	SourceReady string = "SourceReady"
	SourceValid        = "SourceValid"
)

// ModuleSpec defines the desired state of Module
type ModuleSpec struct {
	// +kubebuilder:validation:Optional
	Manifests *Manifests `json:"manifests,omitempty"`
	// +kubebuilder:validation:Optional
	Kustomization *Kustomization `json:"kustomization,omitempty"`
	// +kubebuilder:validation:Optional
	Helm *Helm `json:"helm,omitempty"`
}

type Manifests struct {
	// +kubebuilder:default=false
	Recursive bool `json:"recursive,omitempty"`
	// +kubebuilder:validation:Required
	GitRepository GitRepository `json:"gitRepository"`
}

type Kustomization struct {
	// +kubebuilder:validation:Optional
	Patches *apiextensionsv1.JSON `json:"patches,omitempty"`
	// +kubebuilder:validation:Required
	GitRepository GitRepository `json:"gitRepository"`
}

type Helm struct {
	// +kubebuilder:validation:Optional
	Values *apiextensionsv1.JSON `json:"values,omitempty"`
	// +kubebuilder:validation:Optional
	GitRepository *GitRepository `json:"gitRepository,omitempty"`
	// +kubebuilder:validation:Optional
	HelmRepository *HelmRepository `json:"helmRepository,omitempty"`
}

type GitRepository struct {
	// +kubebuilder:default="60s"
	Interval metav1.Duration `json:"interval,omitempty"`
	// +kubebuilder:validation:Required
	URL string `json:"url"`
	// +kubebuilder:default="/"
	Path string `json:"path,omitempty"`
	// +kubebuilder:validation:Required
	Ref GitRef `json:"ref"`
	// +kubebuilder:validation:Optional
	SecretRef *SecretRef `json:"secretRef,omitempty"`
}

type GitRef struct {
	// +kubebuilder:validation:Enum=branch;commit;tag;semver
	Type string `json:"type"`
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

type HelmRepository struct {
	// +kubebuilder:default="60s"
	Interval metav1.Duration `json:"interval,omitempty"`
	// +kubebuilder:validation:Required
	URL string `json:"url"`
	// +kubebuilder:validation:Optional
	SecretRef *SecretRef `json:"secretRef,omitempty"`
	// +kubebuilder:validation:Required
	HelmChart HelmChart `json:"helmChart"`
}

type HelmChart struct {
	// +kubebuilder:validation:Required
	Chart string `json:"chart"`
	// +kubebuilder:default:=*
	Version string `json:"version,omitempty"`
}

type SecretRef struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// ModuleStatus defines the observed state of Module
type ModuleStatus struct {
	// The phase of a Module is a simple, high-level summary of where the Module is in its lifecycle.
	// +optional
	Phase string `json:"phase,omitempty"`

	// +optional
	Source *Source `json:"source,omitempty"`

	Components []*Component `json:"components,omitempty"`

	// Represents the latest available observations of a Module's current state.
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

type Source struct {
	Path string `json:"path,omitempty"`
}

type Component struct {
	Name       string       `json:"name,omitempty"`
	Containers []*Container `json:"containers,omitempty"`
}

type Container struct {
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:JSONPath=".status.phase",name=Status,type=string

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
