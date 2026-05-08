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

// VerificationKind enumerates the kinds of verification a ResiliencePolicy may declare.
// Each kind maps to a corresponding RecoveryDrill type when the controller
// schedules drills automatically.
// +kubebuilder:validation:Enum=RestoreVerification;FailoverDrill;ReplicaPlacement
type VerificationKind string

const (
	// VerificationRestoreVerification verifies that a backup can be restored into a temporary namespace.
	VerificationRestoreVerification VerificationKind = "RestoreVerification"
	// VerificationFailoverDrill exercises the failover path of a stateful workload.
	VerificationFailoverDrill VerificationKind = "FailoverDrill"
	// VerificationReplicaPlacement verifies that replicas are spread across the configured topology.
	VerificationReplicaPlacement VerificationKind = "ReplicaPlacement"
)

// MitigationKind enumerates corrective actions a ResiliencePolicy may authorize.
// +kubebuilder:validation:Enum=Isolate;ApplyTopologySpread;ScaleSafely;PauseDisruption
type MitigationKind string

const (
	// MitigationIsolate marks a workload as degraded and applies isolation tolerations.
	MitigationIsolate MitigationKind = "Isolate"
	// MitigationApplyTopologySpread applies a topology spread constraint to a workload.
	MitigationApplyTopologySpread MitigationKind = "ApplyTopologySpread"
	// MitigationScaleSafely performs a guarded replica adjustment.
	MitigationScaleSafely MitigationKind = "ScaleSafely"
	// MitigationPauseDisruption pauses voluntary disruption (for example, eviction) for the target.
	MitigationPauseDisruption MitigationKind = "PauseDisruption"
)

// VerificationSpec describes a single verification check the policy authorizes.
type VerificationSpec struct {
	// Kind is the verification kind.
	// +kubebuilder:validation:Required
	Kind VerificationKind `json:"kind"`

	// IntervalSeconds is the minimum interval between executions.
	// +kubebuilder:validation:Minimum=60
	// +kubebuilder:default=3600
	IntervalSeconds int32 `json:"intervalSeconds,omitempty"`

	// FreshnessSeconds is the maximum acceptable age of the last successful execution
	// before the policy reports Degraded.
	// +kubebuilder:validation:Minimum=60
	// +kubebuilder:default=86400
	FreshnessSeconds int32 `json:"freshnessSeconds,omitempty"`
}

// SLOSpec captures the recovery service level objectives the policy enforces.
type SLOSpec struct {
	// RecoveryTimeSeconds is the maximum acceptable time for a recovery action
	// to complete (recovery time objective).
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=900
	RecoveryTimeSeconds int32 `json:"recoveryTimeSeconds,omitempty"`

	// MaxReplicaSkew is the maximum acceptable replica imbalance across the target topology.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=1
	MaxReplicaSkew int32 `json:"maxReplicaSkew,omitempty"`
}

// ResiliencePolicySpec defines the desired state of ResiliencePolicy.
type ResiliencePolicySpec struct {
	// TargetSelector selects the workloads (StatefulSets) governed by this policy.
	// An empty selector matches no workloads. Use the well-known label
	// hilios.io/enabled=true to opt into governance.
	// +kubebuilder:validation:Required
	TargetSelector metav1.LabelSelector `json:"targetSelector"`

	// Verifications are the checks this policy authorises.
	// +kubebuilder:validation:MinItems=1
	Verifications []VerificationSpec `json:"verifications"`

	// Mitigations are the corrective actions this policy authorises.
	// An empty list disables all mitigations (observe-only mode).
	// +optional
	Mitigations []MitigationKind `json:"mitigations,omitempty"`

	// SLO captures recovery service level objectives.
	// +optional
	SLO SLOSpec `json:"slo,omitempty"`

	// Schedule is an optional cron expression that drives periodic re-evaluation.
	// When omitted the policy is evaluated on every reconcile loop.
	// +optional
	Schedule string `json:"schedule,omitempty"`

	// Suspend prevents the controller from creating new RecoveryDrills or
	// applying mitigations. Existing resources are left untouched.
	// +kubebuilder:default=false
	Suspend bool `json:"suspend,omitempty"`
}

// ResiliencePolicyStatus defines the observed state of ResiliencePolicy.
type ResiliencePolicyStatus struct {
	// Conditions represent the latest available observations of the policy state.
	// Standard condition types: Ready, Validated, Degraded.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// ObservedGeneration is the most recent generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// LastEvaluationTime is the last time the controller evaluated the policy.
	// +optional
	LastEvaluationTime *metav1.Time `json:"lastEvaluationTime,omitempty"`

	// MatchedTargets is the number of workloads currently matched by the target selector.
	// +optional
	MatchedTargets int32 `json:"matchedTargets,omitempty"`

	// LastDriftCount is the number of policy violations observed in the last evaluation.
	// +optional
	LastDriftCount int32 `json:"lastDriftCount,omitempty"`

	// LastViolation summarises the most recent policy violation, if any.
	// +optional
	LastViolation string `json:"lastViolation,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=rpol;rpolicies,categories={hilios}
// +kubebuilder:printcolumn:name="Targets",type=integer,JSONPath=`.status.matchedTargets`
// +kubebuilder:printcolumn:name="Drift",type=integer,JSONPath=`.status.lastDriftCount`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// ResiliencePolicy is the Schema for the resiliencepolicies API.
type ResiliencePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResiliencePolicySpec   `json:"spec,omitempty"`
	Status ResiliencePolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ResiliencePolicyList contains a list of ResiliencePolicy.
type ResiliencePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResiliencePolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResiliencePolicy{}, &ResiliencePolicyList{})
}
