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

// Package policy hosts pure helpers used by the ResiliencePolicy controller.
// The freshness check determines, given a verification spec and the most
// recent successful drill timestamp, whether a new drill should be scheduled.
package policy

import (
	"time"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
)

// VerificationStatus pairs a VerificationSpec with the timestamp of the most
// recent successful drill of the same kind. A zero LastSuccess indicates that
// no drill has ever completed for the verification.
type VerificationStatus struct {
	Spec        resiliencev1alpha1.VerificationSpec
	LastSuccess time.Time
}

// IsStale reports whether the verification has not had a successful drill
// within the freshness window. now is supplied so the function is testable
// without relying on the wall clock.
func IsStale(v VerificationStatus, now time.Time) bool {
	if v.Spec.FreshnessSeconds <= 0 {
		return false
	}
	if v.LastSuccess.IsZero() {
		return true
	}
	return now.Sub(v.LastSuccess) > time.Duration(v.Spec.FreshnessSeconds)*time.Second
}

// CountStale returns the number of stale verifications. The result feeds the
// .status.lastDriftCount field on ResiliencePolicy.
func CountStale(items []VerificationStatus, now time.Time) int32 {
	n := int32(0)
	for _, v := range items {
		if IsStale(v, now) {
			n++
		}
	}
	return n
}
