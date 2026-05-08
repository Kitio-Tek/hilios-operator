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

package workload

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ptr(i int32) *int32 { return &i }

func TestSummariseHealthy(t *testing.T) {
	t.Parallel()
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: "db", Namespace: "ns"},
		Spec:       appsv1.StatefulSetSpec{Replicas: ptr(3)},
		Status: appsv1.StatefulSetStatus{
			ReadyReplicas:     3,
			UpdatedReplicas:   3,
			AvailableReplicas: 3,
		},
	}
	s := SummariseStatefulSet(sts)
	if !s.Healthy {
		t.Fatalf("expected healthy summary: %#v", s)
	}
}

func TestSummariseUnhealthy(t *testing.T) {
	t.Parallel()
	sts := &appsv1.StatefulSet{
		Spec: appsv1.StatefulSetSpec{Replicas: ptr(3)},
		Status: appsv1.StatefulSetStatus{
			ReadyReplicas:     2,
			UpdatedReplicas:   3,
			AvailableReplicas: 2,
		},
	}
	if SummariseStatefulSet(sts).Healthy {
		t.Fatal("expected unhealthy summary")
	}
}

func TestPodsForStatefulSet(t *testing.T) {
	t.Parallel()
	sts := &appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: "db"}}
	pods := []corev1.Pod{
		{ObjectMeta: metav1.ObjectMeta{Name: "db-0", OwnerReferences: []metav1.OwnerReference{{Kind: "StatefulSet", Name: "db"}}}},
		{ObjectMeta: metav1.ObjectMeta{Name: "other-0", OwnerReferences: []metav1.OwnerReference{{Kind: "StatefulSet", Name: "other"}}}},
		{ObjectMeta: metav1.ObjectMeta{Name: "stand-alone"}},
	}
	got := PodsForStatefulSet(sts, pods)
	if len(got) != 1 || got[0].Name != "db-0" {
		t.Fatalf("expected only db-0, got %v", got)
	}
}
