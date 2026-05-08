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

package topology

import (
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func nodeWith(name, zone string) corev1.Node {
	return corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: map[string]string{"topology.kubernetes.io/zone": zone},
		},
	}
}

func podOn(name, node string) corev1.Pod {
	return corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       corev1.PodSpec{NodeName: node},
	}
}

func TestDistributeAndSkew(t *testing.T) {
	t.Parallel()

	nodes := []corev1.Node{
		nodeWith("n1", "zone-a"),
		nodeWith("n2", "zone-a"),
		nodeWith("n3", "zone-b"),
	}
	pods := []corev1.Pod{
		podOn("p1", "n1"),
		podOn("p2", "n1"),
		podOn("p3", "n2"),
		podOn("p4", "n3"),
	}

	d := Distribute(pods, nodes, "topology.kubernetes.io/zone")
	if d["zone-a"] != 3 {
		t.Fatalf("zone-a want 3, got %d", d["zone-a"])
	}
	if d["zone-b"] != 1 {
		t.Fatalf("zone-b want 1, got %d", d["zone-b"])
	}
	if got := Skew(d); got != 2 {
		t.Fatalf("skew want 2, got %d", got)
	}
}

func TestSkewEmpty(t *testing.T) {
	t.Parallel()
	if got := Skew(Distribution{}); got != 0 {
		t.Fatalf("empty skew want 0, got %d", got)
	}
}

func TestSortedDomainsDeterministic(t *testing.T) {
	t.Parallel()
	d := Distribution{"c": 1, "a": 0, "b": 2}
	got := SortedDomains(d)
	want := []string{"a", "b", "c"}
	for i, v := range want {
		if got[i] != v {
			t.Fatalf("domain[%d] want %s, got %s", i, v, got[i])
		}
	}
}

func BenchmarkDistribute(b *testing.B) {
	nodes := make([]corev1.Node, 0, 100)
	pods := make([]corev1.Pod, 0, 1000)
	for i := 0; i < 100; i++ {
		zone := "zone-a"
		if i%2 == 0 {
			zone = "zone-b"
		}
		nodes = append(nodes, nodeWith("n", zone))
	}
	for i := 0; i < 1000; i++ {
		pods = append(pods, podOn("p", "n"))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Distribute(pods, nodes, "topology.kubernetes.io/zone")
	}
}

func BenchmarkSkew(b *testing.B) {
	d := Distribution{"a": 5, "b": 7, "c": 3, "d": 9}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Skew(d)
	}
}

func TestDistributeIgnoresUnscheduled(t *testing.T) {
	t.Parallel()
	nodes := []corev1.Node{nodeWith("n1", "zone-a")}
	pods := []corev1.Pod{
		{ObjectMeta: metav1.ObjectMeta{Name: "p1"}}, // no NodeName
		podOn("p2", "n1"),
	}
	d := Distribute(pods, nodes, "topology.kubernetes.io/zone")
	if d["zone-a"] != 1 {
		t.Fatalf("zone-a want 1 (unscheduled excluded), got %d", d["zone-a"])
	}
}

func TestDistributeUsesHostnameDefault(t *testing.T) {
	t.Parallel()
	nodes := []corev1.Node{
		{ObjectMeta: metav1.ObjectMeta{Name: "n1", Labels: map[string]string{"kubernetes.io/hostname": "h1"}}},
	}
	pods := []corev1.Pod{podOn("p1", "n1")}
	d := Distribute(pods, nodes, "")
	if d["h1"] != 1 {
		t.Fatalf("default hostname key should yield h1=1, got %v", d)
	}
}

func TestSkewSingleDomain(t *testing.T) {
	t.Parallel()
	d := Distribution{"only-zone": 5}
	if got := Skew(d); got != 0 {
		t.Fatalf("single-domain skew want 0, got %d", got)
	}
}

func TestSortedDomainsHandlesEmpty(t *testing.T) {
	t.Parallel()
	if got := SortedDomains(Distribution{}); len(got) != 0 {
		t.Fatalf("empty distribution must yield empty slice, got %v", got)
	}
}

func TestSkewWith1DomainsBalanced(t *testing.T) {
	t.Parallel()
	d := Distribution{}
	for j := 0; j < 1; j++ {
		d[fmt.Sprintf("zone-%d", j)] = 3
	}
	if got := Skew(d); got != 0 {
		t.Fatalf("balanced 1 domains must yield zero skew, got %d", got)
	}
}

func TestSkewWith2DomainsBalanced(t *testing.T) {
	t.Parallel()
	d := Distribution{}
	for j := 0; j < 2; j++ {
		d[fmt.Sprintf("zone-%d", j)] = 3
	}
	if got := Skew(d); got != 0 {
		t.Fatalf("balanced 2 domains must yield zero skew, got %d", got)
	}
}

func TestSkewWith3DomainsBalanced(t *testing.T) {
	t.Parallel()
	d := Distribution{}
	for j := 0; j < 3; j++ {
		d[fmt.Sprintf("zone-%d", j)] = 3
	}
	if got := Skew(d); got != 0 {
		t.Fatalf("balanced 3 domains must yield zero skew, got %d", got)
	}
}

func TestSkewWith4DomainsBalanced(t *testing.T) {
	t.Parallel()
	d := Distribution{}
	for j := 0; j < 4; j++ {
		d[fmt.Sprintf("zone-%d", j)] = 3
	}
	if got := Skew(d); got != 0 {
		t.Fatalf("balanced 4 domains must yield zero skew, got %d", got)
	}
}

func TestSkewWith5DomainsBalanced(t *testing.T) {
	t.Parallel()
	d := Distribution{}
	for j := 0; j < 5; j++ {
		d[fmt.Sprintf("zone-%d", j)] = 3
	}
	if got := Skew(d); got != 0 {
		t.Fatalf("balanced 5 domains must yield zero skew, got %d", got)
	}
}
