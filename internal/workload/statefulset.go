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

// Package workload contains read-only helpers that summarise the state of the
// distributed workloads HILIOS governs. These helpers operate on objects
// already loaded into memory so they remain unit-testable without a live API
// server.
package workload

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// HealthSummary captures the key health attributes of a StatefulSet.
type HealthSummary struct {
	Name              string
	Namespace         string
	DesiredReplicas   int32
	ReadyReplicas     int32
	UpdatedReplicas   int32
	AvailableReplicas int32
	Healthy           bool
}

// SummariseStatefulSet returns a HealthSummary built from the StatefulSet's
// observed status. Healthy is true when ready, updated, and available replicas
// equal the desired replica count.
func SummariseStatefulSet(sts *appsv1.StatefulSet) HealthSummary {
	desired := int32(0)
	if sts.Spec.Replicas != nil {
		desired = *sts.Spec.Replicas
	}
	healthy := desired > 0 &&
		sts.Status.ReadyReplicas == desired &&
		sts.Status.UpdatedReplicas == desired &&
		sts.Status.AvailableReplicas == desired
	return HealthSummary{
		Name:              sts.Name,
		Namespace:         sts.Namespace,
		DesiredReplicas:   desired,
		ReadyReplicas:     sts.Status.ReadyReplicas,
		UpdatedReplicas:   sts.Status.UpdatedReplicas,
		AvailableReplicas: sts.Status.AvailableReplicas,
		Healthy:           healthy,
	}
}

// PodsForStatefulSet returns the subset of pods that share the workload's pod
// template label selector. The check is intentionally simple: it matches by
// the well-known StatefulSet labels rather than reproducing the full label
// selector evaluation, since callers already filtered the pod list with the
// selector.
func PodsForStatefulSet(sts *appsv1.StatefulSet, pods []corev1.Pod) []corev1.Pod {
	out := make([]corev1.Pod, 0, len(pods))
	for _, p := range pods {
		for _, owner := range p.OwnerReferences {
			if owner.Kind == "StatefulSet" && owner.Name == sts.Name {
				out = append(out, p)
				break
			}
		}
	}
	return out
}
