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

// Package mitigation contains the small set of corrective actions HILIOS
// authorises ResiliencePolicy to apply. Each action is implemented as a pure
// recommendation (returns a textual description and an optional patch
// fragment) so that controllers can plug recommendations into status without
// taking destructive action by default.
package mitigation

import (
	"fmt"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
)

// Recommendation pairs a mitigation kind with a human readable summary and
// the patch fragment a controller should apply if the policy authorises it.
type Recommendation struct {
	Kind    resiliencev1alpha1.MitigationKind
	Summary string
	Patch   string
}

// Recommend returns the recommendation for the supplied kind and target name.
// An unknown kind returns an empty Recommendation.
func Recommend(kind resiliencev1alpha1.MitigationKind, target string) Recommendation {
	switch kind {
	case resiliencev1alpha1.MitigationApplyTopologySpread:
		return Recommendation{
			Kind:    kind,
			Summary: fmt.Sprintf("apply topologySpreadConstraints to %s with maxSkew=1 across topology.kubernetes.io/zone", target),
			Patch:   topologySpreadPatch(target),
		}
	case resiliencev1alpha1.MitigationIsolate:
		return Recommendation{
			Kind:    kind,
			Summary: fmt.Sprintf("isolate %s from noisy nodes by adding the hilios.io/isolate=true taint toleration", target),
		}
	case resiliencev1alpha1.MitigationScaleSafely:
		return Recommendation{
			Kind:    kind,
			Summary: fmt.Sprintf("scale %s in steps of 1 with PDB-aware drain windows", target),
		}
	case resiliencev1alpha1.MitigationPauseDisruption:
		return Recommendation{
			Kind:    kind,
			Summary: fmt.Sprintf("set a 0-budget PodDisruptionBudget on %s to pause voluntary disruption", target),
		}
	}
	return Recommendation{}
}

// IsAuthorised reports whether the supplied policy authorises the given
// mitigation kind.
func IsAuthorised(p *resiliencev1alpha1.ResiliencePolicy, kind resiliencev1alpha1.MitigationKind) bool {
	for _, m := range p.Spec.Mitigations {
		if m == kind {
			return true
		}
	}
	return false
}

func topologySpreadPatch(target string) string {
	return fmt.Sprintf(`{"spec":{"template":{"spec":{"topologySpreadConstraints":[{"maxSkew":1,"topologyKey":"topology.kubernetes.io/zone","whenUnsatisfiable":"ScheduleAnyway","labelSelector":{"matchLabels":{"app":%q}}}]}}}}`, target)
}
