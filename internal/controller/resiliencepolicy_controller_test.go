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

package controller

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
	"github.com/Kitio-Tek/hilios-operator/internal/conditions"
)

func unitScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	s := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(s); err != nil {
		t.Fatalf("client-go scheme: %v", err)
	}
	if err := resiliencev1alpha1.AddToScheme(s); err != nil {
		t.Fatalf("v1alpha1 scheme: %v", err)
	}
	return s
}

func newPolicy(mut func(*resiliencev1alpha1.ResiliencePolicy)) *resiliencev1alpha1.ResiliencePolicy {
	p := &resiliencev1alpha1.ResiliencePolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "p1", Namespace: "default", Generation: 1},
		Spec: resiliencev1alpha1.ResiliencePolicySpec{
			TargetSelector: metav1.LabelSelector{MatchLabels: map[string]string{"hilios.io/enabled": "true"}},
			Verifications: []resiliencev1alpha1.VerificationSpec{
				{Kind: resiliencev1alpha1.VerificationRestoreVerification, IntervalSeconds: 600, FreshnessSeconds: 86400},
			},
		},
	}
	if mut != nil {
		mut(p)
	}
	return p
}

func TestUnitResiliencePolicyReady(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)

	policy := newPolicy(nil)
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "db",
			Namespace: "default",
			Labels:    map[string]string{"hilios.io/enabled": "true"},
		},
	}

	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(policy, sts).
		WithStatusSubresource(&resiliencev1alpha1.ResiliencePolicy{}).
		Build()

	r := &ResiliencePolicyReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(8)}
	if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: "p1", Namespace: "default"}}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	got := &resiliencev1alpha1.ResiliencePolicy{}
	if err := cli.Get(context.Background(), types.NamespacedName{Name: "p1", Namespace: "default"}, got); err != nil {
		t.Fatalf("get: %v", err)
	}
	if !conditions.IsTrue(got.Status.Conditions, resiliencev1alpha1.ConditionReady) {
		t.Fatalf("Ready condition want True, got %#v", got.Status.Conditions)
	}
	if !conditions.IsTrue(got.Status.Conditions, resiliencev1alpha1.ConditionValidated) {
		t.Fatalf("Validated condition want True, got %#v", got.Status.Conditions)
	}
	if got.Status.MatchedTargets != 1 {
		t.Fatalf("matched targets want 1, got %d", got.Status.MatchedTargets)
	}
}

func TestUnitResiliencePolicyPausedAnnotation(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)
	policy := newPolicy(func(p *resiliencev1alpha1.ResiliencePolicy) {
		p.Annotations = map[string]string{"hilios.io/paused": "true"}
	})
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(policy).
		WithStatusSubresource(&resiliencev1alpha1.ResiliencePolicy{}).
		Build()

	r := &ResiliencePolicyReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(8)}
	if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: "p1", Namespace: "default"}}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	got := &resiliencev1alpha1.ResiliencePolicy{}
	_ = cli.Get(context.Background(), types.NamespacedName{Name: "p1", Namespace: "default"}, got)
	if !conditions.IsFalse(got.Status.Conditions, resiliencev1alpha1.ConditionReady) {
		t.Fatalf("Ready want False (paused annotation), got %#v", got.Status.Conditions)
	}
}

func TestUnitResiliencePolicySuspended(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)
	policy := newPolicy(func(p *resiliencev1alpha1.ResiliencePolicy) {
		p.Spec.Suspend = true
	})
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(policy).
		WithStatusSubresource(&resiliencev1alpha1.ResiliencePolicy{}).
		Build()

	r := &ResiliencePolicyReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(8)}
	if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: "p1", Namespace: "default"}}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	got := &resiliencev1alpha1.ResiliencePolicy{}
	if err := cli.Get(context.Background(), types.NamespacedName{Name: "p1", Namespace: "default"}, got); err != nil {
		t.Fatalf("get: %v", err)
	}
	if !conditions.IsFalse(got.Status.Conditions, resiliencev1alpha1.ConditionReady) {
		t.Fatalf("Ready condition want False (suspended), got %#v", got.Status.Conditions)
	}
}

func TestUnitResiliencePolicyValidationFails(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)
	policy := newPolicy(func(p *resiliencev1alpha1.ResiliencePolicy) {
		p.Spec.Verifications = nil
	})
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(policy).
		WithStatusSubresource(&resiliencev1alpha1.ResiliencePolicy{}).
		Build()

	r := &ResiliencePolicyReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(8)}
	if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: "p1", Namespace: "default"}}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	got := &resiliencev1alpha1.ResiliencePolicy{}
	if err := cli.Get(context.Background(), types.NamespacedName{Name: "p1", Namespace: "default"}, got); err != nil {
		t.Fatalf("get: %v", err)
	}
	if !conditions.IsFalse(got.Status.Conditions, resiliencev1alpha1.ConditionValidated) {
		t.Fatalf("Validated must be False on validation error, got %#v", got.Status.Conditions)
	}
}

func TestUnitResiliencePolicyDriftFromStaleDrills(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)
	policy := newPolicy(func(p *resiliencev1alpha1.ResiliencePolicy) {
		p.Spec.Verifications = []resiliencev1alpha1.VerificationSpec{
			{Kind: resiliencev1alpha1.VerificationRestoreVerification, FreshnessSeconds: 60},
			{Kind: resiliencev1alpha1.VerificationFailoverDrill, FreshnessSeconds: 60},
		}
	})

	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(policy).
		WithStatusSubresource(&resiliencev1alpha1.ResiliencePolicy{}).
		Build()

	r := &ResiliencePolicyReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(8)}
	if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: "p1", Namespace: "default"}}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	got := &resiliencev1alpha1.ResiliencePolicy{}
	_ = cli.Get(context.Background(), types.NamespacedName{Name: "p1", Namespace: "default"}, got)
	if got.Status.LastDriftCount != 2 {
		t.Fatalf("drift count want 2 (no drills, two verifications), got %d", got.Status.LastDriftCount)
	}
	if got.Status.LastViolation == "" {
		t.Fatalf("LastViolation should be populated")
	}
}

func TestUnitResiliencePolicyEmptySelectorDegrades(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)
	policy := newPolicy(func(p *resiliencev1alpha1.ResiliencePolicy) {
		p.Spec.TargetSelector = metav1.LabelSelector{MatchLabels: map[string]string{"hilios.io/enabled": "true"}}
	})
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(policy).
		WithStatusSubresource(&resiliencev1alpha1.ResiliencePolicy{}).
		Build()

	r := &ResiliencePolicyReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(8)}
	if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: "p1", Namespace: "default"}}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	got := &resiliencev1alpha1.ResiliencePolicy{}
	_ = cli.Get(context.Background(), types.NamespacedName{Name: "p1", Namespace: "default"}, got)
	if got.Status.MatchedTargets != 0 {
		t.Fatalf("expected 0 matches, got %d", got.Status.MatchedTargets)
	}
	if !conditions.IsTrue(got.Status.Conditions, resiliencev1alpha1.ConditionDegraded) {
		t.Fatalf("Degraded want True for zero matches, got %#v", got.Status.Conditions)
	}
}
