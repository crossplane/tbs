/*

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	runtimev1alpha1 "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MessageSpec defines the desired state of Message
type MessageSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	runtimev1alpha1.ResourceSpec `json:",inline"`

	// Slack channel that the message should be sent to
	Channel string `json:"channel"`

	// Message text to be sent to the Slack channel
	Text string `json:"text"`
}

// MessageStatus defines the observed state of Message
type MessageStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	runtimev1alpha1.ResourceStatus `json:",inline"`

	// Sent determines whether the message has been sent
	Sent bool `json:"sent,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Message is the Schema for the messages API
type Message struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MessageSpec   `json:"spec,omitempty"`
	Status MessageStatus `json:"status,omitempty"`
}

// SetBindingPhase of this Message.
func (a *Message) SetBindingPhase(p runtimev1alpha1.BindingPhase) {
	a.Status.SetBindingPhase(p)
}

// GetBindingPhase of this Message.
func (a *Message) GetBindingPhase() runtimev1alpha1.BindingPhase {
	return a.Status.GetBindingPhase()
}

// SetConditions of this Message.
func (a *Message) SetConditions(c ...runtimev1alpha1.Condition) {
	a.Status.SetConditions(c...)
}

// SetClaimReference of this Message.
func (a *Message) SetClaimReference(r *corev1.ObjectReference) {
	a.Spec.ClaimReference = r
}

// GetClaimReference of this Message.
func (a *Message) GetClaimReference() *corev1.ObjectReference {
	return a.Spec.ClaimReference
}

// SetClassReference of this Message.
func (a *Message) SetClassReference(r *corev1.ObjectReference) {
	a.Spec.ClassReference = r
}

// GetClassReference of this Message.
func (a *Message) GetClassReference() *corev1.ObjectReference {
	return a.Spec.ClassReference
}

// SetWriteConnectionSecretToReference of this Message.
func (a *Message) SetWriteConnectionSecretToReference(r corev1.LocalObjectReference) {
	a.Spec.WriteConnectionSecretToReference = r
}

// GetWriteConnectionSecretToReference of this Message.
func (a *Message) GetWriteConnectionSecretToReference() corev1.LocalObjectReference {
	return a.Spec.WriteConnectionSecretToReference
}

// GetReclaimPolicy of this Message.
func (a *Message) GetReclaimPolicy() runtimev1alpha1.ReclaimPolicy {
	return a.Spec.ReclaimPolicy
}

// SetReclaimPolicy of this Message.
func (a *Message) SetReclaimPolicy(p runtimev1alpha1.ReclaimPolicy) {
	a.Spec.ReclaimPolicy = p
}

// +kubebuilder:object:root=true

// MessageList contains a list of Message
type MessageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Message `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Message{}, &MessageList{})
}
