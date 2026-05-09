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

func node(name, zone string) *corev1.Node {
	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: map[string]string{"topology.kubernetes.io/zone": zone},
		},
	}
}

func pod(name, nodeName string, lbl map[string]string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default", Labels: lbl},
		Spec:       corev1.PodSpec{NodeName: nodeName},
	}
}

func TestUnitRebalanceCheckBalanced(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)
	check := &resiliencev1alpha1.RebalanceCheck{
		ObjectMeta: metav1.ObjectMeta{Name: "c1", Namespace: "default", Generation: 1},
		Spec: resiliencev1alpha1.RebalanceCheckSpec{
			TargetSelector: metav1.LabelSelector{MatchLabels: map[string]string{"app": "demo"}},
			TopologyKey:    "topology.kubernetes.io/zone",
			MaxSkew:        1,
		},
	}
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(
			check,
			node("n1", "a"), node("n2", "b"),
			pod("p1", "n1", map[string]string{"app": "demo"}),
			pod("p2", "n2", map[string]string{"app": "demo"}),
		).
		WithStatusSubresource(&resiliencev1alpha1.RebalanceCheck{}).
		Build()

	r := &RebalanceCheckReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(8)}
	if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: "c1", Namespace: "default"}}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	got := &resiliencev1alpha1.RebalanceCheck{}
	_ = cli.Get(context.Background(), types.NamespacedName{Name: "c1", Namespace: "default"}, got)
	if !conditions.IsTrue(got.Status.Conditions, resiliencev1alpha1.ConditionBalanced) {
		t.Fatalf("Balanced want True, got %#v", got.Status.Conditions)
	}
	if got.Status.LastSkew != 0 {
		t.Fatalf("skew want 0, got %d", got.Status.LastSkew)
	}
}

func TestUnitRebalanceCheckDrifted(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)
	check := &resiliencev1alpha1.RebalanceCheck{
		ObjectMeta: metav1.ObjectMeta{Name: "c1", Namespace: "default", Generation: 1},
		Spec: resiliencev1alpha1.RebalanceCheckSpec{
			TargetSelector: metav1.LabelSelector{MatchLabels: map[string]string{"app": "demo"}},
			TopologyKey:    "topology.kubernetes.io/zone",
			MaxSkew:        0,
		},
	}
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(
			check,
			node("n1", "a"), node("n2", "b"),
			pod("p1", "n1", map[string]string{"app": "demo"}),
			pod("p2", "n1", map[string]string{"app": "demo"}),
			pod("p3", "n1", map[string]string{"app": "demo"}),
		).
		WithStatusSubresource(&resiliencev1alpha1.RebalanceCheck{}).
		Build()

	r := &RebalanceCheckReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(8)}
	if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: "c1", Namespace: "default"}}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	got := &resiliencev1alpha1.RebalanceCheck{}
	_ = cli.Get(context.Background(), types.NamespacedName{Name: "c1", Namespace: "default"}, got)
	if !conditions.IsFalse(got.Status.Conditions, resiliencev1alpha1.ConditionBalanced) {
		t.Fatalf("Balanced want False, got %#v", got.Status.Conditions)
	}
	if !conditions.IsTrue(got.Status.Conditions, resiliencev1alpha1.ConditionActionRequired) {
		t.Fatalf("ActionRequired want True, got %#v", got.Status.Conditions)
	}
	if got.Status.LastSkew == 0 {
		t.Fatalf("expected non-zero skew")
	}
}

func TestUnitRebalanceCheckEmptySelectorMatchesNothing(t *testing.T) {
	t.Parallel()
	scheme := unitScheme(t)
	check := &resiliencev1alpha1.RebalanceCheck{
		ObjectMeta: metav1.ObjectMeta{Name: "c1", Namespace: "default", Generation: 1},
		Spec: resiliencev1alpha1.RebalanceCheckSpec{
			TargetSelector: metav1.LabelSelector{},
			TopologyKey:    "topology.kubernetes.io/zone",
			MaxSkew:        1,
		},
	}
	cli := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(
			check,
			node("n1", "a"),
			pod("p1", "n1", map[string]string{"app": "demo"}),
		).
		WithStatusSubresource(&resiliencev1alpha1.RebalanceCheck{}).
		Build()

	r := &RebalanceCheckReconciler{Client: cli, Scheme: scheme, Recorder: record.NewFakeRecorder(8)}
	if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: "c1", Namespace: "default"}}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	got := &resiliencev1alpha1.RebalanceCheck{}
	_ = cli.Get(context.Background(), types.NamespacedName{Name: "c1", Namespace: "default"}, got)
	if got.Status.MatchedTargets != 0 {
		t.Fatalf("empty selector must match nothing, got %d targets", got.Status.MatchedTargets)
	}
}
