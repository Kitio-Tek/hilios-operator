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

// DrillValidator implements admission.CustomValidator for RecoveryDrill.
type DrillValidator struct{}

// ValidateCreate runs the shared validation rules.
func (DrillValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	d, ok := obj.(*resiliencev1alpha1.RecoveryDrill)
	if !ok {
		return nil, errors.New("expected a RecoveryDrill")
	}
	return drillValidate(d)
}

// ValidateUpdate refuses changes to spec.type once the drill has started.
func (DrillValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	oldD, ok := oldObj.(*resiliencev1alpha1.RecoveryDrill)
	if !ok {
		return nil, errors.New("old object is not a RecoveryDrill")
	}
	newD, ok := newObj.(*resiliencev1alpha1.RecoveryDrill)
	if !ok {
		return nil, errors.New("new object is not a RecoveryDrill")
	}
	if oldD.Status.Phase != "" && oldD.Status.Phase != resiliencev1alpha1.DrillPhasePending && oldD.Spec.Type != newD.Spec.Type {
		return nil, errors.New("spec.type is immutable once the drill has started")
	}
	return drillValidate(newD)
}

// ValidateDelete is a no-op.
func (DrillValidator) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func drillValidate(d *resiliencev1alpha1.RecoveryDrill) (admission.Warnings, error) {
	errs := validation.DrillSpec(d)
	var warnings admission.Warnings
	fatal := make([]string, 0, len(errs))
	for _, e := range errs {
		msg := e.Error()
		if strings.HasPrefix(msg, "spec.source is recommended") {
			warnings = append(warnings, msg)
			continue
		}
		fatal = append(fatal, msg)
	}
	if len(fatal) == 0 {
		return warnings, nil
	}
	return warnings, errors.New(strings.Join(fatal, "; "))
}
