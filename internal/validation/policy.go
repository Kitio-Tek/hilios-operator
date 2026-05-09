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

// Package validation hosts reusable spec validation routines shared between
// controllers and (when enabled) the admission webhook. Each function reports
// errors as a flat slice so callers can present every problem at once instead
// of failing on the first issue.
package validation

import (
	"fmt"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
)

// PolicySpec validates a ResiliencePolicy spec and returns every detected
// problem. An empty result indicates a valid spec.
func PolicySpec(p *resiliencev1alpha1.ResiliencePolicy) []error {
	var errs []error
	if p == nil {
		return []error{fmt.Errorf("policy is nil")}
	}
	if len(p.Spec.Verifications) == 0 {
		errs = append(errs, fmt.Errorf("spec.verifications must declare at least one verification"))
	}
	for i, v := range p.Spec.Verifications {
		if v.Kind == "" {
			errs = append(errs, fmt.Errorf("spec.verifications[%d].kind is required", i))
		}
		if v.IntervalSeconds < 0 {
			errs = append(errs, fmt.Errorf("spec.verifications[%d].intervalSeconds must be non-negative", i))
		}
		if v.FreshnessSeconds < 0 {
			errs = append(errs, fmt.Errorf("spec.verifications[%d].freshnessSeconds must be non-negative", i))
		}
	}
	if p.Spec.SLO.RecoveryTimeSeconds < 0 {
		errs = append(errs, fmt.Errorf("spec.slo.recoveryTimeSeconds must be non-negative"))
	}
	if p.Spec.SLO.MaxReplicaSkew < 0 {
		errs = append(errs, fmt.Errorf("spec.slo.maxReplicaSkew must be non-negative"))
	}
	return errs
}

// DrillSpec validates a RecoveryDrill spec and returns every detected problem.
func DrillSpec(d *resiliencev1alpha1.RecoveryDrill) []error {
	var errs []error
	if d == nil {
		return []error{fmt.Errorf("drill is nil")}
	}
	if d.Spec.Type == "" {
		errs = append(errs, fmt.Errorf("spec.type is required"))
	}
	if d.Spec.Type == resiliencev1alpha1.DrillRestoreVerification && d.Spec.Source == nil {
		// RestoreVerification without an explicit source is allowed for
		// templates, but the controller will mark it as Skipped at execution
		// time. Emit an informational error so callers can decide whether to
		// treat templates as valid.
		errs = append(errs, fmt.Errorf("spec.source is recommended for RestoreVerification drills"))
	}
	if d.Spec.TimeoutSeconds < 0 {
		errs = append(errs, fmt.Errorf("spec.timeoutSeconds must be non-negative"))
	}
	for i, hc := range d.Spec.HealthChecks {
		switch hc.Type {
		case resiliencev1alpha1.ProbeTypeHTTP:
			if hc.URL == "" {
				errs = append(errs, fmt.Errorf("spec.healthChecks[%d].url is required for HTTP probes", i))
			}
		case resiliencev1alpha1.ProbeTypeKubernetes:
			if hc.Resource == nil {
				errs = append(errs, fmt.Errorf("spec.healthChecks[%d].resource is required for Kubernetes probes", i))
			}
		case resiliencev1alpha1.ProbeTypeCmd:
			if hc.Command == "" {
				errs = append(errs, fmt.Errorf("spec.healthChecks[%d].command is required for Cmd probes", i))
			}
		case resiliencev1alpha1.ProbeTypePod:
			if hc.PodSelector == nil {
				errs = append(errs, fmt.Errorf("spec.healthChecks[%d].podSelector is required for Pod probes", i))
			}
		case "":
			errs = append(errs, fmt.Errorf("spec.healthChecks[%d].type is required", i))
		}
	}
	return errs
}
