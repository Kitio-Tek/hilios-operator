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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
	"github.com/Kitio-Tek/hilios-operator/internal/conditions"
)

func contendedPod(name, reason string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default", Labels: map[string]string{"app": "demo"}},
		Spec:       corev1.PodSpec{NodeName: "n1"},
		Status: corev1.PodStatus{
			Conditions: []corev1.PodCondition{
				{Type: corev1.ContainersReady, Status: corev1.ConditionFalse, Reason: reason, Message: "test"},
			},
		},
	}
}

func TestUnitContentionReportNoFindings(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)
	rep := &resiliencev1alpha1.ContentionReport{
		ObjectMeta: metav1.ObjectMeta{Name: "r1", Namespace: "default", Generation: 1},
		Spec: resiliencev1alpha1.ContentionReportSpec{
			TargetSelector: metav1.LabelSelector{MatchLabels: map[string]string{"app": "demo"}},
			RecommendOnly:  true,
		},
	}
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(rep, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "p1", Namespace: "default", Labels: map[string]string{"app": "demo"}}}).
		WithStatusSubresource(&resiliencev1alpha1.ContentionReport{}).
		Build()

	r := &ContentionReportReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(8)}
	if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: "r1", Namespace: "default"}}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	got := &resiliencev1alpha1.ContentionReport{}
	_ = cli.Get(context.Background(), types.NamespacedName{Name: "r1", Namespace: "default"}, got)
	if !conditions.IsFalse(got.Status.Conditions, resiliencev1alpha1.ConditionObserved) {
		t.Fatalf("Observed want False (no findings), got %#v", got.Status.Conditions)
	}
}

func TestUnitContentionReportFindings(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)
	rep := &resiliencev1alpha1.ContentionReport{
		ObjectMeta: metav1.ObjectMeta{Name: "r1", Namespace: "default", Generation: 1},
		Spec: resiliencev1alpha1.ContentionReportSpec{
			TargetSelector: metav1.LabelSelector{MatchLabels: map[string]string{"app": "demo"}},
			RecommendOnly:  true,
		},
	}
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(rep, contendedPod("p1", "Throttled"), contendedPod("p2", "MemoryPressure")).
		WithStatusSubresource(&resiliencev1alpha1.ContentionReport{}).
		Build()

	r := &ContentionReportReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(8)}
	if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: "r1", Namespace: "default"}}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	got := &resiliencev1alpha1.ContentionReport{}
	_ = cli.Get(context.Background(), types.NamespacedName{Name: "r1", Namespace: "default"}, got)
	if !conditions.IsTrue(got.Status.Conditions, resiliencev1alpha1.ConditionObserved) {
		t.Fatalf("Observed want True, got %#v", got.Status.Conditions)
	}
	if len(got.Status.Findings) != 2 {
		t.Fatalf("findings want 2, got %d", len(got.Status.Findings))
	}
	if got.Status.Findings[0].Recommendation == "" {
		t.Fatalf("expected recommendation populated")
	}
	// Mitigated must be False because HILIOS does not actually apply mitigations.
	if !conditions.IsFalse(got.Status.Conditions, resiliencev1alpha1.ConditionMitigated) {
		t.Fatalf("Mitigated want False (HILIOS records recommendations only), got %#v", got.Status.Conditions)
	}
}

func TestUnitContentionReportMitigatedFalseWhenNotRecommendOnly(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)
	rep := &resiliencev1alpha1.ContentionReport{
		ObjectMeta: metav1.ObjectMeta{Name: "r1", Namespace: "default", Generation: 1},
		Spec: resiliencev1alpha1.ContentionReportSpec{
			TargetSelector: metav1.LabelSelector{MatchLabels: map[string]string{"app": "demo"}},
			RecommendOnly:  false,
		},
	}
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(rep, contendedPod("p1", "Throttled")).
		WithStatusSubresource(&resiliencev1alpha1.ContentionReport{}).
		Build()

	r := &ContentionReportReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(8)}
	if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: "r1", Namespace: "default"}}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	got := &resiliencev1alpha1.ContentionReport{}
	_ = cli.Get(context.Background(), types.NamespacedName{Name: "r1", Namespace: "default"}, got)
	if !conditions.IsFalse(got.Status.Conditions, resiliencev1alpha1.ConditionMitigated) {
		t.Fatalf("Mitigated must be False even when RecommendOnly=false (mitigation not yet implemented), got %#v", got.Status.Conditions)
	}
}
