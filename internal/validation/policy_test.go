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

package validation

import (
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
)

func TestPolicySpecValid(t *testing.T) {
	t.Parallel()
	p := &resiliencev1alpha1.ResiliencePolicy{
		Spec: resiliencev1alpha1.ResiliencePolicySpec{
			TargetSelector: metav1.LabelSelector{MatchLabels: map[string]string{"app": "x"}},
			Verifications: []resiliencev1alpha1.VerificationSpec{
				{Kind: resiliencev1alpha1.VerificationRestoreVerification, IntervalSeconds: 60, FreshnessSeconds: 600},
			},
			SLO: resiliencev1alpha1.SLOSpec{RecoveryTimeSeconds: 60, MaxReplicaSkew: 1},
		},
	}
	if errs := PolicySpec(p); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestPolicySpecMissingVerifications(t *testing.T) {
	t.Parallel()
	p := &resiliencev1alpha1.ResiliencePolicy{}
	errs := PolicySpec(p)
	if len(errs) == 0 {
		t.Fatal("expected at least one error")
	}
}

func TestPolicySpecNil(t *testing.T) {
	t.Parallel()
	if errs := PolicySpec(nil); len(errs) != 1 {
		t.Fatalf("expected 1 error for nil, got %v", errs)
	}
}

func TestDrillSpecValid(t *testing.T) {
	t.Parallel()
	d := &resiliencev1alpha1.RecoveryDrill{
		Spec: resiliencev1alpha1.RecoveryDrillSpec{
			Type: resiliencev1alpha1.DrillFailoverDrill,
			HealthChecks: []resiliencev1alpha1.HealthCheck{
				{Name: "h", Type: "HTTP", URL: "http://x"},
			},
			TimeoutSeconds: 60,
		},
	}
	if errs := DrillSpec(d); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestDrillSpecBadProbes(t *testing.T) {
	t.Parallel()
	d := &resiliencev1alpha1.RecoveryDrill{
		Spec: resiliencev1alpha1.RecoveryDrillSpec{
			Type: resiliencev1alpha1.DrillFailoverDrill,
			HealthChecks: []resiliencev1alpha1.HealthCheck{
				{Name: "no-url", Type: "HTTP"},
				{Name: "no-resource", Type: "Kubernetes"},
				{Name: "no-cmd", Type: "Cmd"},
				{Name: "no-selector", Type: "Pod"},
				{Name: "no-type"},
			},
			TimeoutSeconds: 60,
		},
	}
	errs := DrillSpec(d)
	if len(errs) != 5 {
		t.Fatalf("expected 5 errors, got %d: %v", len(errs), errs)
	}
}

func TestPolicySpecNegativeFreshness(t *testing.T) {
	t.Parallel()
	p := &resiliencev1alpha1.ResiliencePolicy{
		Spec: resiliencev1alpha1.ResiliencePolicySpec{
			Verifications: []resiliencev1alpha1.VerificationSpec{
				{Kind: resiliencev1alpha1.VerificationRestoreVerification, FreshnessSeconds: -1},
			},
		},
	}
	if errs := PolicySpec(p); len(errs) == 0 {
		t.Fatal("negative freshness must be rejected")
	}
}

func TestDrillSpecRejectsEmptyType(t *testing.T) {
	t.Parallel()
	d := &resiliencev1alpha1.RecoveryDrill{}
	errs := DrillSpec(d)
	found := false
	for _, e := range errs {
		if strings.Contains(e.Error(), "spec.type is required") {
			found = true
		}
	}
	if !found {
		t.Fatal("missing type error not reported")
	}
}
