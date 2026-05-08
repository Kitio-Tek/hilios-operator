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

package healthcheck

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
)

func TestRunKubernetesFound(t *testing.T) {
	t.Parallel()
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatalf("scheme: %v", err)
	}
	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "p", Namespace: "ns"}}
	cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(pod).Build()

	hc := resiliencev1alpha1.HealthCheck{
		Name: "exists", Type: "Kubernetes",
		Resource: &resiliencev1alpha1.KubernetesResourceRef{
			APIVersion: "v1", Kind: "Pod", Namespace: "ns", Name: "p",
		},
	}
	if err := RunKubernetes(context.Background(), cli, hc); err != nil {
		t.Fatalf("expected found, got %v", err)
	}
}

func TestRunKubernetesNotFound(t *testing.T) {
	t.Parallel()
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	cli := fake.NewClientBuilder().WithScheme(scheme).Build()
	hc := resiliencev1alpha1.HealthCheck{
		Name: "missing", Type: "Kubernetes",
		Resource: &resiliencev1alpha1.KubernetesResourceRef{
			APIVersion: "v1", Kind: "Pod", Namespace: "ns", Name: "absent",
		},
	}
	if err := RunKubernetes(context.Background(), cli, hc); err == nil {
		t.Fatal("expected not-found error")
	}
}

func TestRunKubernetesValidation(t *testing.T) {
	t.Parallel()
	cases := []resiliencev1alpha1.HealthCheck{
		{Name: "wrong-type", Type: "HTTP"},
		{Name: "missing-resource", Type: "Kubernetes"},
		{Name: "bad-apiversion", Type: "Kubernetes", Resource: &resiliencev1alpha1.KubernetesResourceRef{APIVersion: "a/b/c", Kind: "Pod", Name: "p"}},
		{Name: "missing-kind", Type: "Kubernetes", Resource: &resiliencev1alpha1.KubernetesResourceRef{APIVersion: "v1", Name: "p"}},
	}
	for _, hc := range cases {
		t.Run(hc.Name, func(t *testing.T) {
			if err := RunKubernetes(context.Background(), nil, hc); err == nil {
				t.Fatalf("expected validation error for %s", hc.Name)
			}
		})
	}
}

func TestRunKubernetesEmptyAPIVersion(t *testing.T) {
	t.Parallel()
	hc := resiliencev1alpha1.HealthCheck{
		Name: "x", Type: "Kubernetes",
		Resource: &resiliencev1alpha1.KubernetesResourceRef{Kind: "Pod", Name: "p"},
	}
	if err := RunKubernetes(context.Background(), nil, hc); err == nil {
		t.Fatal("empty apiVersion must return an error")
	}
}
