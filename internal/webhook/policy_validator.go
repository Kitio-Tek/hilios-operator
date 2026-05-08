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

// Package webhook provides validating admission webhook handlers for HILIOS
// resources. The handlers reuse the spec validation logic in
// internal/validation so the controller and the webhook always agree about
// what counts as a valid spec.
//
// The handlers are intentionally not yet wired in cmd/main.go: the project
// does not require admission webhooks for v1alpha1, and shipping them
// without certificate machinery would only increase Helm chart complexity.
// Once the API graduates to v1, register the handlers via
// ctrl.NewWebhookManagedBy(mgr).For(...).WithValidator(...).
package webhook

import (
	"context"
	"errors"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
	"github.com/Kitio-Tek/hilios-operator/internal/validation"
)

// PolicyValidator implements admission.CustomValidator for ResiliencePolicy.
type PolicyValidator struct{}

// ValidateCreate runs the shared validation rules.
func (PolicyValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	p, ok := obj.(*resiliencev1alpha1.ResiliencePolicy)
	if !ok {
		return nil, errors.New("expected a ResiliencePolicy")
	}
	return policyValidate(p)
}

// ValidateUpdate runs the shared validation rules and refuses changes to immutable fields.
func (PolicyValidator) ValidateUpdate(_ context.Context, _ runtime.Object, newObj runtime.Object) (admission.Warnings, error) {
	p, ok := newObj.(*resiliencev1alpha1.ResiliencePolicy)
	if !ok {
		return nil, errors.New("expected a ResiliencePolicy")
	}
	return policyValidate(p)
}

// ValidateDelete is a no-op.
func (PolicyValidator) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func policyValidate(p *resiliencev1alpha1.ResiliencePolicy) (admission.Warnings, error) {
	errs := validation.PolicySpec(p)
	if len(errs) == 0 {
		return nil, nil
	}
	msgs := make([]string, 0, len(errs))
	for _, e := range errs {
		msgs = append(msgs, e.Error())
	}
	return nil, errors.New(strings.Join(msgs, "; "))
}
