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
	"testing"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
)

func TestDrillValidateCreateValid(t *testing.T) {
	t.Parallel()
	d := &resiliencev1alpha1.RecoveryDrill{
		Spec: resiliencev1alpha1.RecoveryDrillSpec{
			Type:           resiliencev1alpha1.DrillFailoverDrill,
			TimeoutSeconds: 60,
		},
	}
	if _, err := (DrillValidator{}).ValidateCreate(context.Background(), d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDrillValidateUpdateImmutableType(t *testing.T) {
	t.Parallel()
	oldD := &resiliencev1alpha1.RecoveryDrill{
		Spec:   resiliencev1alpha1.RecoveryDrillSpec{Type: resiliencev1alpha1.DrillRestoreVerification, TimeoutSeconds: 60},
		Status: resiliencev1alpha1.RecoveryDrillStatus{Phase: resiliencev1alpha1.DrillPhaseRunning},
	}
	newD := oldD.DeepCopy()
	newD.Spec.Type = resiliencev1alpha1.DrillFailoverDrill
	if _, err := (DrillValidator{}).ValidateUpdate(context.Background(), oldD, newD); err == nil {
		t.Fatal("expected immutability error")
	}
}

func TestDrillValidateRestoreSourceWarning(t *testing.T) {
	t.Parallel()
	d := &resiliencev1alpha1.RecoveryDrill{
		Spec: resiliencev1alpha1.RecoveryDrillSpec{
			Type:           resiliencev1alpha1.DrillRestoreVerification,
			TimeoutSeconds: 60,
		},
	}
	warnings, err := (DrillValidator{}).ValidateCreate(context.Background(), d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) == 0 {
		t.Fatal("expected a warning about missing source")
	}
}
