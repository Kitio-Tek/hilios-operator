/*
Copyright 2026.

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
)

// TopologyDistribution captures the per-domain replica count observed during a check.
type TopologyDistribution struct {
	// Domain is the value of the topology key (for example, a zone or hostname).
	Domain string `json:"domain"`
	// Replicas is the number of replicas observed in the domain.
	Replicas int32 `json:"replicas"`
}

// RebalanceCheckSpec defines the desired state of RebalanceCheck.
type RebalanceCheckSpec struct {
	// TargetSelector selects the workloads inspected by this check.
	// +kubebuilder:validation:Required
	TargetSelector metav1.LabelSelector `json:"targetSelector"`

	// TopologyKey is the node label key used to compute the distribution.
	// +kubebuilder:default="kubernetes.io/hostname"
	TopologyKey string `json:"topologyKey,omitempty"`

	// MaxSkew is the maximum acceptable replica imbalance across topology domains.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=1
	MaxSkew int32 `json:"maxSkew,omitempty"`

	// Schedule is an optional cron expression that drives periodic checks.
	// +optional
	Schedule string `json:"schedule,omitempty"`

	// DryRun reports drift but never recommends or applies a mitigation.
	// +kubebuilder:default=true
	DryRun bool `json:"dryRun,omitempty"`
}

// RebalanceCheckStatus defines the observed state of RebalanceCheck.
type RebalanceCheckStatus struct {
	// Conditions track the lifecycle of the check.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// ObservedGeneration is the most recent generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// LastEvaluationTime is the most recent time the controller evaluated the check.
	// +optional
	LastEvaluationTime *metav1.Time `json:"lastEvaluationTime,omitempty"`

	// LastSkew is the most recently observed skew across topology domains.
	// +optional
	LastSkew int32 `json:"lastSkew,omitempty"`

	// MatchedTargets is the number of workloads matched by the target selector.
	// +optional
	MatchedTargets int32 `json:"matchedTargets,omitempty"`

	// Distribution captures the most recent per-domain replica counts.
	// +optional
	Distribution []TopologyDistribution `json:"distribution,omitempty"`

	// Message is a human readable summary of the latest evaluation.
	// +optional
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=rcheck;rchecks,categories={hilios}
// +kubebuilder:printcolumn:name="Targets",type=integer,JSONPath=`.status.matchedTargets`
// +kubebuilder:printcolumn:name="Skew",type=integer,JSONPath=`.status.lastSkew`
// +kubebuilder:printcolumn:name="Balanced",type=string,JSONPath=`.status.conditions[?(@.type=="Balanced")].status`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// RebalanceCheck is the Schema for the rebalancechecks API.
type RebalanceCheck struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RebalanceCheckSpec   `json:"spec,omitempty"`
	Status RebalanceCheckStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RebalanceCheckList contains a list of RebalanceCheck.
type RebalanceCheckList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RebalanceCheck `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RebalanceCheck{}, &RebalanceCheckList{})
}
