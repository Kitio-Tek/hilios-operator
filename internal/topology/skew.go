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

// Package topology computes replica distribution metrics across topology
// domains. These helpers are intentionally pure (no Kubernetes API access) so
// that they can be unit tested without a fake client or envtest harness.
package topology

import (
	"sort"

	corev1 "k8s.io/api/core/v1"
)

// Distribution is a domain -> replica count mapping. Domains with zero
// replicas are kept in the map so that the absence of pods in a zone counts
// toward the skew computation.
type Distribution map[string]int32

// Distribute computes how many of the supplied pods land in each value of the
// node label whose key is topologyKey. Pods scheduled on nodes that are not in
// nodes are bucketed under the empty domain "".
func Distribute(pods []corev1.Pod, nodes []corev1.Node, topologyKey string) Distribution {
	if topologyKey == "" {
		topologyKey = corev1.LabelHostname
	}
	domains := make(Distribution, len(nodes))
	nodeDomain := make(map[string]string, len(nodes))
	for _, n := range nodes {
		nodeDomain[n.Name] = n.Labels[topologyKey]
		domains[n.Labels[topologyKey]] = 0
	}
	for _, p := range pods {
		if p.Spec.NodeName == "" {
			continue
		}
		domain := nodeDomain[p.Spec.NodeName]
		domains[domain]++
	}
	return domains
}

// Skew returns the worst observed imbalance across topology domains.
// It is the difference between the most populated and the least
// populated domain. An empty distribution returns 0. Domains with zero replicas
// are included so an underutilised zone counts toward the skew.
func Skew(d Distribution) int32 {
	if len(d) == 0 {
		return 0
	}
	var min, max int32 = -1, -1
	for _, v := range d {
		if min == -1 || v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	if min < 0 {
		min = 0
	}
	return max - min
}

// SortedDomains returns the distribution as a slice ordered by domain name to
// produce deterministic output for status fields.
func SortedDomains(d Distribution) []string {
	out := make([]string, 0, len(d))
	for k := range d {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
