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

// Package contention exposes a small interface that the ContentionReport
// reconciler delegates pod-level evaluation to. The default implementation
// reads PodConditions; alternative implementations can wrap a metrics backend
// without changing the controller.
package contention

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
)

// Contention reasons surfaced on ContentionFinding.Reason. They mirror the
// PodCondition reasons that container runtimes commonly set when a workload
// is throttled, under memory pressure, or hosted on a noisy node.
const (
	ReasonThrottled      = "Throttled"
	ReasonMemoryPressure = "MemoryPressure"
	ReasonCPUSteal       = "CPUSteal"
)

// Evaluator returns a finding when the pod is experiencing contention or nil
// otherwise. now is supplied so callers can record deterministic timestamps.
type Evaluator interface {
	Evaluate(pod *corev1.Pod, now time.Time) *resiliencev1alpha1.ContentionFinding
}

// EvaluatorFunc adapts a function to the Evaluator interface.
type EvaluatorFunc func(pod *corev1.Pod, now time.Time) *resiliencev1alpha1.ContentionFinding

// Evaluate satisfies the Evaluator interface.
func (f EvaluatorFunc) Evaluate(pod *corev1.Pod, now time.Time) *resiliencev1alpha1.ContentionFinding {
	return f(pod, now)
}

// PodConditionEvaluator is the default implementation. It inspects
// ContainersReady=False with a known reason.
type PodConditionEvaluator struct{}

// Evaluate inspects pod.Status.Conditions for known contention reasons.
func (PodConditionEvaluator) Evaluate(pod *corev1.Pod, now time.Time) *resiliencev1alpha1.ContentionFinding {
	for _, c := range pod.Status.Conditions {
		if c.Type != corev1.ContainersReady {
			continue
		}
		if c.Status == corev1.ConditionTrue {
			continue
		}
		switch c.Reason {
		case ReasonThrottled, ReasonMemoryPressure, ReasonCPUSteal:
			return &resiliencev1alpha1.ContentionFinding{
				Pod:            pod.Name,
				Node:           pod.Spec.NodeName,
				Reason:         c.Reason,
				ObservedAt:     metav1.NewTime(now),
				Recommendation: Recommendation(c.Reason),
				Value:          c.Message,
			}
		}
	}
	return nil
}

// Recommendation returns the canonical operator hint for a given contention reason.
func Recommendation(reason string) string {
	switch reason {
	case ReasonThrottled:
		return "increase CPU limits or relax CPU manager static pinning"
	case ReasonMemoryPressure:
		return "raise memory requests or move pod to a less crowded node"
	case ReasonCPUSteal:
		return "isolate to a node with stable steal time or apply pod priority"
	}
	return ""
}
