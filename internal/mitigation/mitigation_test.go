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

package mitigation

import (
	"strings"
	"testing"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
)

func TestRecommendKnown(t *testing.T) {
	t.Parallel()
	cases := []resiliencev1alpha1.MitigationKind{
		resiliencev1alpha1.MitigationApplyTopologySpread,
		resiliencev1alpha1.MitigationIsolate,
		resiliencev1alpha1.MitigationScaleSafely,
		resiliencev1alpha1.MitigationPauseDisruption,
	}
	for _, k := range cases {
		t.Run(string(k), func(t *testing.T) {
			r := Recommend(k, "payments")
			if r.Kind != k {
				t.Fatalf("kind want %s, got %s", k, r.Kind)
			}
			if r.Summary == "" || !strings.Contains(r.Summary, "payments") {
				t.Fatalf("summary missing target: %q", r.Summary)
			}
		})
	}
}

func TestRecommendUnknown(t *testing.T) {
	t.Parallel()
	r := Recommend("not-a-kind", "x")
	if r.Kind != "" {
		t.Fatalf("unknown must return empty: %#v", r)
	}
}

func TestIsAuthorised(t *testing.T) {
	t.Parallel()
	p := &resiliencev1alpha1.ResiliencePolicy{
		Spec: resiliencev1alpha1.ResiliencePolicySpec{
			Mitigations: []resiliencev1alpha1.MitigationKind{resiliencev1alpha1.MitigationApplyTopologySpread},
		},
	}
	if !IsAuthorised(p, resiliencev1alpha1.MitigationApplyTopologySpread) {
		t.Fatal("authorised mitigation should return true")
	}
	if IsAuthorised(p, resiliencev1alpha1.MitigationIsolate) {
		t.Fatal("unauthorised mitigation should return false")
	}
}

func TestRecommendIncludesPatchForTopologySpread(t *testing.T) {
	t.Parallel()
	r := Recommend(resiliencev1alpha1.MitigationApplyTopologySpread, "demo")
	if r.Patch == "" {
		t.Fatal("ApplyTopologySpread must produce a patch fragment")
	}
}
