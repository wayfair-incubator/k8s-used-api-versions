/*
Copyright 2022.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// UsedApiVersionsSpec defines the desired state of UsedApiVersions
type UsedApiVersionsSpec struct {

	// UsedApiVersions is a list of API versions
	UsedApiVersions []APIVersionMeta `json:"usedApiVersions,omitempty"`
}

// APIVersionMeta defines the used API version and Kind
type APIVersionMeta struct {

	// APIVersion is the name of the API version used by specific kind.
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind is the Object type such as "Deployment" or "Ingress"
	Kind string `json:"kind,omitempty"`
}

// UsedApiVersionsStatus defines the observed state of UsedApiVersions
type UsedApiVersionsStatus struct {
	// ApiVersionsStatus is a list of checked API versions
	ApiVersionsStatus []APIVersionStatus `json:"apiVersionsStatus,omitempty"`
	// FinalStatus is the overall status for all the used API versions
	FinalStatus FinalStatusResult `json:"finalStatus,omitempty"`
}

// FinalStatusResult is the overall status for all the used API versions
type FinalStatusResult struct {
	// Number of deprecated API Versions
	Deprecated int `json:"deprecated" yaml:"deprecated"`
	// Number of removed API Versions
	Removed int `json:"removed" yaml:"removed"`
	// Number of removed API Versions in the next release
	RemovedInNextRelease int `json:"removedInNextRelease" yaml:"removedInNextRelease"`
	// Number of removed API Versions in the next two releases
	RemovedInNextTwoReleases int `json:"removedInNextTwoReleases" yaml:"removedInNextTwoReleases"`
}

// APIVersionStatus defines the observed API version status
type APIVersionStatus struct {
	// APIVersion is the name of the apiVersion.
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	// Kind is the Object type
	Kind string `json:"kind" yaml:"kind"`
	// Whether the API Version is deprecated or not
	Deprecated bool `json:"deprecated" yaml:"deprecated"`
	// Whether the API Version is removed or not
	Removed bool `json:"removed" yaml:"removed"`
	// Kubernetes version in which the API is deprecated in
	DeprecatedInVersion string `json:"deprecatedInVersion" yaml:"deprecatedInVersion"`
	// Kubernetes version in which the API is removed in
	RemovedInVersion string `json:"removedInVersion" yaml:"removedInVersion"`
	// ReplacementAPI is the new supported apiVersion.
	ReplacementAPI string `json:"replacementApi" yaml:"replacementApi"`
	// Whether the apiVersion will be removed in the next release or not
	RemovedInNextRelease bool `json:"removedInNextRelease" yaml:"removedInNextRelease"`
	// Whether the apiVersion will be removed in the next release or not
	RemovedInNextTwoReleases bool `json:"removedInNextTwoReleases" yaml:"removedInNextTwoReleases"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:resource:shortName=uav
// +kubebuilder:printcolumn:name="Kind",type=string,JSONPath=`.kind`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="Deprecated",type=integer,JSONPath=`.status.finalStatus.deprecated`
// +kubebuilder:printcolumn:name="Removed",type=integer,JSONPath=`.status.finalStatus.removed`
// +kubebuilder:printcolumn:name="Removed-NEXT-Release",type=integer,JSONPath=`.status.finalStatus.removedInNextRelease`,priority=10
// +kubebuilder:printcolumn:name="Removed-NEXT-Two-Releases",type=integer,JSONPath=`.status.finalStatus.removedInNextTwoReleases`,priority=10
type UsedApiVersions struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UsedApiVersionsSpec   `json:"spec,omitempty"`
	Status UsedApiVersionsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// UsedApiVersionsList contains a list of UsedApiVersions
type UsedApiVersionsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UsedApiVersions `json:"items"`
}

func init() {
	SchemeBuilder.Register(&UsedApiVersions{}, &UsedApiVersionsList{})
}
