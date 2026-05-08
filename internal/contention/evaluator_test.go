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

package contention

import (
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPodConditionEvaluatorThrottled(t *testing.T) {
	t.Parallel()
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "p1"},
		Spec:       corev1.PodSpec{NodeName: "n1"},
		Status: corev1.PodStatus{
			Conditions: []corev1.PodCondition{
				{Type: corev1.ContainersReady, Status: corev1.ConditionFalse, Reason: "Throttled", Message: "cpu"},
			},
		},
	}
	got := PodConditionEvaluator{}.Evaluate(pod, time.Now())
	if got == nil {
		t.Fatal("expected finding for Throttled")
	}
	if got.Recommendation == "" {
		t.Fatal("recommendation should be populated")
	}
}

func TestPodConditionEvaluatorReady(t *testing.T) {
	t.Parallel()
	pod := &corev1.Pod{
		Status: corev1.PodStatus{
			Conditions: []corev1.PodCondition{
				{Type: corev1.ContainersReady, Status: corev1.ConditionTrue},
			},
		},
	}
	if got := (PodConditionEvaluator{}).Evaluate(pod, time.Now()); got != nil {
		t.Fatalf("expected no finding for ready pod, got %#v", got)
	}
}

func TestRecommendationKnown(t *testing.T) {
	t.Parallel()
	cases := []string{"Throttled", "MemoryPressure", "CPUSteal"}
	for _, c := range cases {
		if Recommendation(c) == "" {
			t.Fatalf("missing recommendation for %s", c)
		}
	}
	if Recommendation("OOMKilled") != "" {
		t.Fatal("unknown reason must return empty recommendation")
	}
}
