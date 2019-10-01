/*
Copyright 2019 The Crossplane Authors.

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
	"reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	runtimev1alpha1 "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

var (
	InstanceKind             = reflect.TypeOf(Instance{}).Name()
	InstanceKindAPIVersion   = InstanceKind + "." + GroupVersion.String()
	InstanceGroupVersionKind = GroupVersion.WithKind(InstanceKind)
)

// +kubebuilder:validation:Required

type InstanceParameters struct {
	// Label is the unique name of this Linode Instance
	// +optional
	Label string `json:"label,omitempty"`

	// Image is the disk image to be applied to the first instance disk
	// +optional
	Image string `json:"image,omitempty"`

	// AuthorizedUsers are Linode user accounts whose SSH keys will be authorized to SSH into the instance
	// +optional
	AuthorizedUsers []string `json:"authorizedUsers,omitempty"`

	// Region defines the geographic location of a Linode Instance
	Region string `json:"region"`

	// Type is the Linode Instance Type which represents the cost, processor, memory, transfer, and storage profile of the Instance
	Type string `json:"type"`

	// Status is the current activity status of a Linode Instance
	// +kubebuilder:validation:Enum=offline;running
	Status string `json:"status,omitempty"`
}

// InstanceSpec defines the desired state of Instance
type InstanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	runtimev1alpha1.ResourceSpec `json:",inline"`
	InstanceParameters           `json:",inline"`
}

// InstanceStatus defines the observed state of Instance
type InstanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	runtimev1alpha1.ResourceStatus `json:",inline"`

	// Id is the unique immutable numeric identifier of a Linode Instance
	// +optional
	Id int `json:"id,omitempty"`

	// Status is the current activity status of a Linode Instance
	Status string `json:"status"`

	// Label is the unique mutable name of a Linode Instance
	Label string `json:"label"`

	// Region defines the geographic location of a Linode Instance
	Region string `json:"region"`

	// Type is the Linode Instance Type which represents the cost, processor, memory, transfer, and storage profile of the Instance
	Type string `json:"type"`

	// IPv6 is the public IPv6 address of a Linode Instance
	// +optional
	IPv6 string `json:"ipv6,omitempty"`

	// IPv4 is the list of IPv4 addresses associated with a Linode Instance
	// +optional
	IPv4 []string `json:"ipv4,omitempty"`

	// Image is the image detected on a Linode Instance disk
	// +optional
	Image string `json:"image,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="LABEL",type="string",JSONPath=".status.label",description="Unique label associated with this Linode Instance",priority=1
// +kubebuilder:printcolumn:name="REGION",type="string",JSONPath=".status.region",description="Region where this Linode Instance is deployed",priority=1
// +kubebuilder:printcolumn:name="IPV4",type="string",JSONPath=".status.ipv4[0]",description="First IPv4 address of this Linode Instance",priority=1
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.status",description="Power status of this Linode Instance",priority=1

// Instance is the Schema for the instances API
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec InstanceSpec `json:"spec,omitempty"`

	// +optional
	Status InstanceStatus `json:"status,omitempty"`
}

// GetSpec returns the MySQL server's spec.
func (s *Instance) GetSpec() *InstanceSpec {
	return &s.Spec
}

// GetStatus returns the MySQL server's status.
func (s *Instance) GetStatus() *InstanceStatus {
	return &s.Status
}

// SetStatus sets the MySQL server's status.
func (s *Instance) SetStatus(status *InstanceStatus) {
	s.Status = *status
}

// SetBindingPhase of this Instance.
func (a *Instance) SetBindingPhase(p runtimev1alpha1.BindingPhase) {
	a.Status.SetBindingPhase(p)
}

// GetBindingPhase of this Instance.
func (a *Instance) GetBindingPhase() runtimev1alpha1.BindingPhase {
	return a.Status.GetBindingPhase()
}

// SetConditions of this Instance.
func (a *Instance) SetConditions(c ...runtimev1alpha1.Condition) {
	a.Status.SetConditions(c...)
}

// SetClaimReference of this Instance.
func (a *Instance) SetClaimReference(r *corev1.ObjectReference) {
	a.Spec.ClaimReference = r
}

// GetClaimReference of this Instance.
func (a *Instance) GetClaimReference() *corev1.ObjectReference {
	return a.Spec.ClaimReference
}

// SetNonPortableClassReference of this Instance.
func (a *Instance) SetNonPortableClassReference(r *corev1.ObjectReference) {
	a.Spec.NonPortableClassReference = r
}

// GetNonPortableClassReference of this Instance.
func (a *Instance) GetNonPortableClassReference() *corev1.ObjectReference {
	return a.Spec.NonPortableClassReference
}

// SetWriteConnectionSecretToReference of this Instance.
func (a *Instance) SetWriteConnectionSecretToReference(r corev1.LocalObjectReference) {
	a.Spec.WriteConnectionSecretToReference = r
}

// GetWriteConnectionSecretToReference of this Instance.
func (a *Instance) GetWriteConnectionSecretToReference() corev1.LocalObjectReference {
	return a.Spec.WriteConnectionSecretToReference
}

// GetReclaimPolicy of this Instance.
func (a *Instance) GetReclaimPolicy() runtimev1alpha1.ReclaimPolicy {
	return a.Spec.ReclaimPolicy
}

// SetReclaimPolicy of this Instance.
func (a *Instance) SetReclaimPolicy(p runtimev1alpha1.ReclaimPolicy) {
	a.Spec.ReclaimPolicy = p
}

// +kubebuilder:object:root=true

// InstanceList contains a list of Instance
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Instance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Instance{}, &InstanceList{})
}
