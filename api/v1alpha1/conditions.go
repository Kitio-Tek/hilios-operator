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

// Standard condition types reported by HILIOS resources. Centralising the
// type strings here lets both controllers and tests use a single source of
// truth so misspellings cannot slip into status.
const (
	// ConditionReady indicates the resource has been processed at least once
	// and the controller is tracking it.
	ConditionReady = "Ready"

	// ConditionValidated indicates the resource passes spec validation.
	ConditionValidated = "Validated"

	// ConditionDegraded indicates the resource has detected a policy violation
	// or an unhealthy state that requires operator attention.
	ConditionDegraded = "Degraded"

	// ConditionScheduled indicates a workflow object has been scheduled.
	ConditionScheduled = "Scheduled"

	// ConditionRunning indicates a workflow object is currently executing.
	ConditionRunning = "Running"

	// ConditionSucceeded indicates a workflow completed successfully.
	ConditionSucceeded = "Succeeded"

	// ConditionFailed indicates a workflow failed.
	ConditionFailed = "Failed"

	// ConditionBalanced indicates the topology of a workload meets the policy.
	ConditionBalanced = "Balanced"

	// ConditionDrifted indicates the topology of a workload no longer meets the policy.
	ConditionDrifted = "Drifted"

	// ConditionActionRequired indicates a corrective action is recommended or required.
	ConditionActionRequired = "ActionRequired"

	// ConditionObserved indicates contention has been observed on the workload.
	ConditionObserved = "Observed"

	// ConditionMitigated indicates contention has been mitigated.
	ConditionMitigated = "Mitigated"

	// ConditionEscalated indicates contention persisted past the policy threshold and was escalated.
	ConditionEscalated = "Escalated"
)

// Standard condition reasons.
const (
	ReasonReconciling          = "Reconciling"
	ReasonReady                = "Ready"
	ReasonValidationFailed     = "ValidationFailed"
	ReasonSelectorEmpty        = "SelectorEmpty"
	ReasonSuspended            = "Suspended"
	ReasonRestoreFailed        = "RestoreFailed"
	ReasonRestoreVerified      = "RestoreVerified"
	ReasonReplicaSkew          = "ReplicaSkewDetected"
	ReasonReplicaBalanced      = "ReplicasBalanced"
	ReasonContentionDetected   = "ContentionDetected"
	ReasonMitigationApplied    = "MitigationApplied"
	ReasonMitigationDisallowed = "MitigationDisallowed"
	ReasonEscalated            = "Escalated"
	ReasonTimeoutExceeded      = "TimeoutExceeded"
	ReasonHealthCheckFailed    = "HealthCheckFailed"
	ReasonScheduled            = "Scheduled"
	ReasonStarted              = "Started"
	ReasonCompleted            = "Completed"
)
