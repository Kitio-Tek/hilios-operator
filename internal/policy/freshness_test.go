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

package policy

import (
	"testing"
	"time"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
)

func TestIsStaleZeroFreshnessNeverStale(t *testing.T) {
	t.Parallel()
	v := VerificationStatus{Spec: resiliencev1alpha1.VerificationSpec{FreshnessSeconds: 0}}
	if IsStale(v, time.Now()) {
		t.Fatal("zero freshness must mean never stale")
	}
}

func TestIsStaleNoSuccessIsStale(t *testing.T) {
	t.Parallel()
	v := VerificationStatus{Spec: resiliencev1alpha1.VerificationSpec{FreshnessSeconds: 60}}
	if !IsStale(v, time.Now()) {
		t.Fatal("no LastSuccess and positive freshness must be stale")
	}
}

func TestIsStaleWithinWindow(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	v := VerificationStatus{
		Spec:        resiliencev1alpha1.VerificationSpec{FreshnessSeconds: 3600},
		LastSuccess: now.Add(-30 * time.Minute),
	}
	if IsStale(v, now) {
		t.Fatal("within window must not be stale")
	}
}

func TestIsStaleOutsideWindow(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	v := VerificationStatus{
		Spec:        resiliencev1alpha1.VerificationSpec{FreshnessSeconds: 3600},
		LastSuccess: now.Add(-2 * time.Hour),
	}
	if !IsStale(v, now) {
		t.Fatal("outside window must be stale")
	}
}

func TestCountStale(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	items := []VerificationStatus{
		{Spec: resiliencev1alpha1.VerificationSpec{FreshnessSeconds: 3600}, LastSuccess: now.Add(-30 * time.Minute)},
		{Spec: resiliencev1alpha1.VerificationSpec{FreshnessSeconds: 3600}, LastSuccess: now.Add(-2 * time.Hour)},
		{Spec: resiliencev1alpha1.VerificationSpec{FreshnessSeconds: 0}},
	}
	if got := CountStale(items, now); got != 1 {
		t.Fatalf("CountStale want 1, got %d", got)
	}
}
