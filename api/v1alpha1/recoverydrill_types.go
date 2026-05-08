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

// DrillType selects the kind of corrective verification a RecoveryDrill performs.
// +kubebuilder:validation:Enum=RestoreVerification;FailoverDrill
type DrillType string

const (
	// DrillRestoreVerification verifies that a backup can be restored into a temporary namespace.
	DrillRestoreVerification DrillType = "RestoreVerification"
	// DrillFailoverDrill verifies that the failover path of a stateful workload completes within the SLO.
	DrillFailoverDrill DrillType = "FailoverDrill"
)

// DrillPhase is the high-level state of a drill.
// +kubebuilder:validation:Enum=Pending;Running;Succeeded;Failed
type DrillPhase string

const (
	// DrillPhasePending indicates the drill has been created but not yet started.
	DrillPhasePending DrillPhase = "Pending"
	// DrillPhaseRunning indicates the drill is currently executing.
	DrillPhaseRunning DrillPhase = "Running"
	// DrillPhaseSucceeded indicates the drill completed successfully and produced evidence.
	DrillPhaseSucceeded DrillPhase = "Succeeded"
	// DrillPhaseFailed indicates the drill did not complete successfully.
	DrillPhaseFailed DrillPhase = "Failed"
)

// BackupSource references a backup to be restored as part of the drill.
type BackupSource struct {
	// Name of the backup object (for example, a Velero Backup).
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Namespace of the backup object.
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`

	// APIGroup of the backup object. Defaults to "velero.io".
	// +kubebuilder:default="velero.io"
	APIGroup string `json:"apiGroup,omitempty"`

	// Kind of the backup object. Defaults to "Backup".
	// +kubebuilder:default="Backup"
	Kind string `json:"kind,omitempty"`
}

// HealthCheck describes a probe executed against the restored workload.
// The probe model is inspired by the Litmus Chaos probe schema (httpProbe,
// k8sProbe, cmdProbe) so users familiar with chaos engineering pipelines can
// reuse mental models. Unlike Litmus, HILIOS probes are read-only and never
// inject failures.
type HealthCheck struct {
	// Name is a human readable identifier for the health check.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Type selects the probe implementation.
	// +kubebuilder:validation:Enum=HTTP;Pod;Cmd;Kubernetes
	// +kubebuilder:default=Pod
	Type string `json:"type"`

	// URL is the HTTP URL to probe (only for type HTTP).
	// +optional
	URL string `json:"url,omitempty"`

	// ExpectedStatusCode is the HTTP status code that indicates success.
	// +kubebuilder:default=200
	// +optional
	ExpectedStatusCode int32 `json:"expectedStatusCode,omitempty"`

	// PodSelector selects pods to probe (only for type Pod).
	// +optional
	PodSelector *metav1.LabelSelector `json:"podSelector,omitempty"`

	// Command is the shell command executed inside the probe runner pod
	// (only for type Cmd). Commands run in a busybox image without elevated
	// privileges and must complete within TimeoutSeconds.
	// +optional
	Command string `json:"command,omitempty"`

	// ExpectedOutput is a substring matched against the command's stdout
	// (only for type Cmd). When empty, exit code 0 is sufficient.
	// +optional
	ExpectedOutput string `json:"expectedOutput,omitempty"`

	// Resource targets a Kubernetes object whose existence is verified
	// (only for type Kubernetes).
	// +optional
	Resource *KubernetesResourceRef `json:"resource,omitempty"`

	// TimeoutSeconds bounds an individual probe attempt.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=30
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`

	// RetryCount is the number of times the probe is re-attempted on failure
	// before being marked failed.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=0
	RetryCount int32 `json:"retryCount,omitempty"`
}

// KubernetesResourceRef identifies a target object for the Kubernetes probe.
type KubernetesResourceRef struct {
	// APIVersion of the target object, for example "apps/v1".
	// +kubebuilder:validation:Required
	APIVersion string `json:"apiVersion"`
	// Kind of the target object, for example "StatefulSet".
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`
	// Namespace of the target object.
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// Name of the target object.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// EvidenceRecord is an audit-friendly artefact produced by a drill.
type EvidenceRecord struct {
	// Step is the drill step that produced the evidence.
	Step string `json:"step"`
	// Time the evidence was produced.
	Time metav1.Time `json:"time"`
	// Result captures the outcome (Pass, Fail, Skip).
	Result string `json:"result"`
	// Message is a human readable summary.
	// +optional
	Message string `json:"message,omitempty"`
}

// LocalObjectReference points to a local object by name.
type LocalObjectReference struct {
	// Name of the referenced object.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// RecoveryDrillSpec defines the desired state of RecoveryDrill.
type RecoveryDrillSpec struct {
	// Type is the kind of drill.
	// +kubebuilder:validation:Required
	Type DrillType `json:"type"`

	// PolicyRef optionally links the drill to a ResiliencePolicy.
	// +optional
	PolicyRef *LocalObjectReference `json:"policyRef,omitempty"`

	// Source describes the backup to restore (RestoreVerification only).
	// +optional
	Source *BackupSource `json:"source,omitempty"`

	// VerificationNamespace is the namespace where verification artefacts are created.
	// When empty, the controller derives a unique name from the drill metadata.
	// +optional
	VerificationNamespace string `json:"verificationNamespace,omitempty"`

	// HealthChecks are probes evaluated against the restored or failed-over workload.
	// +optional
	HealthChecks []HealthCheck `json:"healthChecks,omitempty"`

	// Cleanup controls whether verification artefacts are deleted after the drill completes.
	// +kubebuilder:default=true
	Cleanup bool `json:"cleanup,omitempty"`

	// TimeoutSeconds bounds the entire drill.
	// +kubebuilder:validation:Minimum=60
	// +kubebuilder:default=1800
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`

	// Schedule is an optional cron expression for periodic drills.
	// +optional
	Schedule string `json:"schedule,omitempty"`
}

// RecoveryDrillStatus defines the observed state of RecoveryDrill.
type RecoveryDrillStatus struct {
	// Phase is the high-level state of the drill.
	// +optional
	Phase DrillPhase `json:"phase,omitempty"`

	// Conditions track the lifecycle of the drill.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// ObservedGeneration is the most recent generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// StartTime is when the drill began executing.
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// CompletionTime is when the drill reached a terminal phase.
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// DurationSeconds is the wall clock duration of the drill.
	// +optional
	DurationSeconds int32 `json:"durationSeconds,omitempty"`

	// Evidence is the ordered list of audit records produced by the drill.
	// +optional
	Evidence []EvidenceRecord `json:"evidence,omitempty"`

	// Message is a human readable summary of the most recent state transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=rdrill;rdrills,categories={hilios}
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.type`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Duration",type=integer,JSONPath=`.status.durationSeconds`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// RecoveryDrill is the Schema for the recoverydrills API.
type RecoveryDrill struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RecoveryDrillSpec   `json:"spec,omitempty"`
	Status RecoveryDrillStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RecoveryDrillList contains a list of RecoveryDrill.
type RecoveryDrillList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RecoveryDrill `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RecoveryDrill{}, &RecoveryDrillList{})
}
