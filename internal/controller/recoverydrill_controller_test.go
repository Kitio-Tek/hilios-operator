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

func newDrill(name string, mut func(*resiliencev1alpha1.RecoveryDrill)) *resiliencev1alpha1.RecoveryDrill {
	d := &resiliencev1alpha1.RecoveryDrill{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default", Generation: 1},
		Spec: resiliencev1alpha1.RecoveryDrillSpec{
			Type:           resiliencev1alpha1.DrillRestoreVerification,
			Cleanup:        true,
			TimeoutSeconds: 600,
		},
	}
	if mut != nil {
		mut(d)
	}
	return d
}

func reconcileUntilDone(t *testing.T, r *RecoveryDrillReconciler, name string) *resiliencev1alpha1.RecoveryDrill {
	t.Helper()
	for i := 0; i < 10; i++ {
		if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: "default"}}); err != nil {
			t.Fatalf("reconcile #%d: %v", i, err)
		}
		got := &resiliencev1alpha1.RecoveryDrill{}
		if err := r.Get(context.Background(), types.NamespacedName{Name: name, Namespace: "default"}, got); err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.Status.Phase == resiliencev1alpha1.DrillPhaseSucceeded || got.Status.Phase == resiliencev1alpha1.DrillPhaseFailed {
			return got
		}
	}
	t.Fatalf("drill did not reach terminal phase")
	return nil
}

func TestUnitRecoveryDrillSucceeds(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)
	drill := newDrill("d1", nil)
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(drill).
		WithStatusSubresource(&resiliencev1alpha1.RecoveryDrill{}).
		Build()

	r := &RecoveryDrillReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(16)}
	got := reconcileUntilDone(t, r, "d1")

	if got.Status.Phase != resiliencev1alpha1.DrillPhaseSucceeded {
		t.Fatalf("phase want Succeeded, got %s", got.Status.Phase)
	}
	if !conditions.IsTrue(got.Status.Conditions, resiliencev1alpha1.ConditionSucceeded) {
		t.Fatalf("Succeeded condition not True: %#v", got.Status.Conditions)
	}
	if got.Status.DurationSeconds < 0 {
		t.Fatalf("duration should be non-negative")
	}
	if len(got.Status.Evidence) < 2 {
		t.Fatalf("expected at least 2 evidence records, got %d", len(got.Status.Evidence))
	}
}

func TestUnitRecoveryDrillFailsOnHealthCheck(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)
	drill := newDrill("d1", func(d *resiliencev1alpha1.RecoveryDrill) {
		d.Spec.HealthChecks = []resiliencev1alpha1.HealthCheck{
			{Name: "missing", Type: "Kubernetes", Resource: &resiliencev1alpha1.KubernetesResourceRef{
				APIVersion: "v1", Kind: "Namespace", Name: "definitely-not-here",
			}},
		}
	})
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(drill).
		WithStatusSubresource(&resiliencev1alpha1.RecoveryDrill{}).
		Build()

	r := &RecoveryDrillReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(16)}
	got := reconcileUntilDone(t, r, "d1")

	if got.Status.Phase != resiliencev1alpha1.DrillPhaseFailed {
		t.Fatalf("phase want Failed, got %s", got.Status.Phase)
	}
	if !conditions.IsTrue(got.Status.Conditions, resiliencev1alpha1.ConditionFailed) {
		t.Fatalf("Failed condition not True: %#v", got.Status.Conditions)
	}
	foundFail := false
	for _, ev := range got.Status.Evidence {
		if ev.Result == "Fail" {
			foundFail = true
		}
	}
	if !foundFail {
		t.Fatalf("expected at least one Fail evidence record, got %#v", got.Status.Evidence)
	}
}

func TestUnitRecoveryDrillCreatesVerificationNamespace(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)
	drill := newDrill("d1", func(d *resiliencev1alpha1.RecoveryDrill) {
		d.Spec.VerificationNamespace = "verify-ns"
	})
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(drill).
		WithStatusSubresource(&resiliencev1alpha1.RecoveryDrill{}).
		Build()

	r := &RecoveryDrillReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(16)}
	_ = reconcileUntilDone(t, r, "d1")

	ns := &corev1.Namespace{}
	if err := cli.Get(context.Background(), types.NamespacedName{Name: "verify-ns"}, ns); err != nil {
		t.Fatalf("verification namespace not created: %v", err)
	}
	if ns.Labels["app.kubernetes.io/managed-by"] != "hilios-operator" {
		t.Fatalf("verification namespace missing managed-by label: %#v", ns.Labels)
	}
}
