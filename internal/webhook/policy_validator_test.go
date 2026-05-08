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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
)

func TestValidateCreateRejectsWrongType(t *testing.T) {
	t.Parallel()
	if _, err := (PolicyValidator{}).ValidateCreate(context.Background(), &corev1.ConfigMap{}); err == nil {
		t.Fatal("expected error when validating non-ResiliencePolicy object")
	}
}

func TestValidateCreateValid(t *testing.T) {
	t.Parallel()
	p := &resiliencev1alpha1.ResiliencePolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "p"},
		Spec: resiliencev1alpha1.ResiliencePolicySpec{
			TargetSelector: metav1.LabelSelector{MatchLabels: map[string]string{"a": "b"}},
			Verifications: []resiliencev1alpha1.VerificationSpec{
				{Kind: resiliencev1alpha1.VerificationRestoreVerification, IntervalSeconds: 60, FreshnessSeconds: 600},
			},
		},
	}
	if _, err := (PolicyValidator{}).ValidateCreate(context.Background(), p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCreateInvalid(t *testing.T) {
	t.Parallel()
	p := &resiliencev1alpha1.ResiliencePolicy{ObjectMeta: metav1.ObjectMeta{Name: "p"}}
	if _, err := (PolicyValidator{}).ValidateCreate(context.Background(), p); err == nil {
		t.Fatal("expected error for empty spec")
	}
}

func TestValidateUpdate(t *testing.T) {
	t.Parallel()
	old := &resiliencev1alpha1.ResiliencePolicy{}
	cur := &resiliencev1alpha1.ResiliencePolicy{}
	if _, err := (PolicyValidator{}).ValidateUpdate(context.Background(), old, cur); err == nil {
		t.Fatal("expected error for empty spec on update")
	}
}

func TestValidateDeleteNoOp(t *testing.T) {
	t.Parallel()
	if _, err := (PolicyValidator{}).ValidateDelete(context.Background(), nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
