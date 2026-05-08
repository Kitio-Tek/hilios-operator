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

// Package labels declares the well-known labels and annotations that HILIOS
// uses to opt resources into governance, attribute managed objects, and signal
// pause states. Centralising these strings keeps controllers free of magic
// constants.
package labels

const (
	// LabelEnabled is the opt-in label. Workloads that carry the value "true"
	// are eligible for ResiliencePolicy matching and corrective workflows.
	// HILIOS never matches workloads that do not opt in: the absence of the
	// label is interpreted as "this workload is out of scope".
	LabelEnabled = "hilios.io/enabled"

	// LabelManagedBy identifies the controller responsible for an object.
	// It uses the standard Kubernetes recommended label key so external tools
	// (Lens, k9s, ArgoCD) can group HILIOS-managed objects automatically.
	LabelManagedBy = "app.kubernetes.io/managed-by"

	// LabelManagedByValue is the canonical value of LabelManagedBy.
	LabelManagedByValue = "hilios-operator"

	// LabelPolicy attributes a managed object to its parent ResiliencePolicy.
	LabelPolicy = "hilios.io/policy"

	// LabelDrill attributes a managed object to its parent RecoveryDrill.
	LabelDrill = "hilios.io/drill"

	// LabelCheck attributes a managed object to its parent RebalanceCheck.
	LabelCheck = "hilios.io/check"

	// LabelReport attributes a managed object to its parent ContentionReport.
	LabelReport = "hilios.io/report"

	// AnnotationPaused, when set to "true" on a HILIOS resource, suspends
	// reconciliation without removing the object.
	AnnotationPaused = "hilios.io/paused"

	// AnnotationVerificationNamespace records the temporary namespace created
	// for a restore verification drill.
	AnnotationVerificationNamespace = "hilios.io/verification-namespace"
)

// ManagedBy returns the canonical labels used to mark an object as managed by
// HILIOS. Callers are expected to merge these into existing label maps.
func ManagedBy() map[string]string {
	return map[string]string{
		LabelManagedBy: LabelManagedByValue,
	}
}

// MergeManagedBy returns a copy of base with the managed-by label applied.
func MergeManagedBy(base map[string]string) map[string]string {
	out := make(map[string]string, len(base)+1)
	for k, v := range base {
		out[k] = v
	}
	out[LabelManagedBy] = LabelManagedByValue
	return out
}
