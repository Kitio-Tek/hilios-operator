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

// ContentionThresholds describes when a workload is considered to be experiencing contention.
type ContentionThresholds struct {
	// CPUStealPercent is the threshold for CPU steal time, expressed as a percentage.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=10
	CPUStealPercent int32 `json:"cpuStealPercent,omitempty"`

	// MemoryPressureEvents is the threshold for the number of memory pressure events
	// observed within the evaluation window.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=5
	MemoryPressureEvents int32 `json:"memoryPressureEvents,omitempty"`

	// ThrottlingPercent is the threshold for the percentage of CPU time
	// spent throttled by the kernel.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=15
	ThrottlingPercent int32 `json:"throttlingPercent,omitempty"`
}

// ContentionFinding records a single observation of contention.
type ContentionFinding struct {
	// Pod is the name of the impacted pod.
	Pod string `json:"pod"`
	// Node is the name of the node hosting the pod.
	Node string `json:"node"`
	// Reason captures the contention signal (CPUSteal, ThrottledCPU, MemoryPressure).
	Reason string `json:"reason"`
	// Value is the observed numeric value associated with the reason.
	// +optional
	Value string `json:"value,omitempty"`
	// ObservedAt is when the signal was observed.
	ObservedAt metav1.Time `json:"observedAt"`
	// Recommendation is a short human readable mitigation suggestion.
	// +optional
	Recommendation string `json:"recommendation,omitempty"`
}

// ContentionReportSpec defines the desired state of ContentionReport.
type ContentionReportSpec struct {
	// TargetSelector selects pods to evaluate.
	// +kubebuilder:validation:Required
	TargetSelector metav1.LabelSelector `json:"targetSelector"`

	// WindowMinutes is the look-back window used when evaluating signals.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=15
	WindowMinutes int32 `json:"windowMinutes,omitempty"`

	// Thresholds describe when a workload is considered to be experiencing contention.
	// +optional
	Thresholds ContentionThresholds `json:"thresholds,omitempty"`

	// RecommendOnly disables any active mitigation and limits the controller to
	// recording recommendations on findings.
	// +kubebuilder:default=true
	RecommendOnly bool `json:"recommendOnly,omitempty"`

	// Schedule is an optional cron expression for periodic re-evaluation.
	// +optional
	Schedule string `json:"schedule,omitempty"`
}

// ContentionReportStatus defines the observed state of ContentionReport.
type ContentionReportStatus struct {
	// Conditions track the lifecycle of the report.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// ObservedGeneration is the most recent generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// LastEvaluationTime is the most recent time the controller evaluated the report.
	// +optional
	LastEvaluationTime *metav1.Time `json:"lastEvaluationTime,omitempty"`

	// MatchedTargets is the number of pods matched by the selector.
	// +optional
	MatchedTargets int32 `json:"matchedTargets,omitempty"`

	// Findings is the list of contention findings observed during the most recent evaluation.
	// +optional
	Findings []ContentionFinding `json:"findings,omitempty"`

	// Message is a human readable summary of the most recent evaluation.
	// +optional
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=creport;creports,categories={hilios}
// +kubebuilder:printcolumn:name="Targets",type=integer,JSONPath=`.status.matchedTargets`
// +kubebuilder:printcolumn:name="Findings",type=integer,JSONPath=`.status.findings[*].pod`,priority=1
// +kubebuilder:printcolumn:name="Observed",type=string,JSONPath=`.status.conditions[?(@.type=="Observed")].status`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// ContentionReport is the Schema for the contentionreports API.
type ContentionReport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContentionReportSpec   `json:"spec,omitempty"`
	Status ContentionReportStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ContentionReportList contains a list of ContentionReport.
type ContentionReportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContentionReport `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContentionReport{}, &ContentionReportList{})
}
