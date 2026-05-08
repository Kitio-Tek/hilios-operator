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

// Package conditions contains helpers for managing the metav1.Condition slices
// that HILIOS resources expose on their .status. Helpers in this package are
// safe for concurrent use only when callers serialise access to the underlying
// slice (typical reconciler pattern).
package conditions

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Set installs or updates a condition on the supplied slice. The slice is
// returned for fluent assignment.
func Set(conds *[]metav1.Condition, condType, reason, message string, status metav1.ConditionStatus, observedGen int64) []metav1.Condition {
	if conds == nil {
		return nil
	}
	cond := metav1.Condition{
		Type:               condType,
		Status:             status,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: observedGen,
	}
	meta.SetStatusCondition(conds, cond)
	return *conds
}

// True is a convenience for Set with ConditionStatus True.
func True(conds *[]metav1.Condition, condType, reason, message string, observedGen int64) []metav1.Condition {
	return Set(conds, condType, reason, message, metav1.ConditionTrue, observedGen)
}

// False is a convenience for Set with ConditionStatus False.
func False(conds *[]metav1.Condition, condType, reason, message string, observedGen int64) []metav1.Condition {
	return Set(conds, condType, reason, message, metav1.ConditionFalse, observedGen)
}

// Unknown is a convenience for Set with ConditionStatus Unknown.
func Unknown(conds *[]metav1.Condition, condType, reason, message string, observedGen int64) []metav1.Condition {
	return Set(conds, condType, reason, message, metav1.ConditionUnknown, observedGen)
}

// IsTrue reports whether a condition of the given type is present and True.
func IsTrue(conds []metav1.Condition, condType string) bool {
	return meta.IsStatusConditionTrue(conds, condType)
}

// IsFalse reports whether a condition of the given type is present and False.
func IsFalse(conds []metav1.Condition, condType string) bool {
	return meta.IsStatusConditionFalse(conds, condType)
}

// Find returns a copy of the named condition or nil if not present.
func Find(conds []metav1.Condition, condType string) *metav1.Condition {
	return meta.FindStatusCondition(conds, condType)
}

// Remove deletes any condition with the given type from the slice in place
// and returns the updated slice.
func Remove(conds []metav1.Condition, condType string) []metav1.Condition {
	out := conds[:0]
	for _, c := range conds {
		if c.Type == condType {
			continue
		}
		out = append(out, c)
	}
	return out
}
