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

package selector

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func TestBuildEmptyMatchesNothing(t *testing.T) {
	t.Parallel()
	s, err := Build(metav1.LabelSelector{})
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if s.Matches(labels.Set{"app": "x"}) {
		t.Fatal("empty selector must match nothing")
	}
}

func TestBuildMatchLabels(t *testing.T) {
	t.Parallel()
	s, err := Build(metav1.LabelSelector{MatchLabels: map[string]string{"app": "x"}})
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if !s.Matches(labels.Set{"app": "x"}) {
		t.Fatal("matching set must match")
	}
	if s.Matches(labels.Set{"app": "y"}) {
		t.Fatal("non-matching set must not match")
	}
}

func TestBuildMatchExpressions(t *testing.T) {
	t.Parallel()
	in := metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{
		{Key: "tier", Operator: metav1.LabelSelectorOpIn, Values: []string{"critical"}},
	}}
	s, err := Build(in)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if !s.Matches(labels.Set{"tier": "critical"}) {
		t.Fatal("In matches critical")
	}
	if s.Matches(labels.Set{"tier": "low"}) {
		t.Fatal("In must not match low")
	}
}

func TestIsEmpty(t *testing.T) {
	t.Parallel()
	if !IsEmpty(metav1.LabelSelector{}) {
		t.Fatal("empty selector must be empty")
	}
	if IsEmpty(metav1.LabelSelector{MatchLabels: map[string]string{"a": "b"}}) {
		t.Fatal("populated selector must not be empty")
	}
}
