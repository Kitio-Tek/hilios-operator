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

// Package selector wraps the conversion from metav1.LabelSelector to a
// labels.Selector with HILIOS-specific defaults: an empty selector matches
// nothing (instead of the labels package default of matching everything),
// because reconcilers iterate over the result and accidentally matching every
// resource in the cluster is rarely what the operator wants.
package selector

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// Build converts a metav1.LabelSelector into a labels.Selector. An empty input
// (no MatchLabels and no MatchExpressions) returns a selector that matches
// nothing.
func Build(in metav1.LabelSelector) (labels.Selector, error) {
	if len(in.MatchLabels) == 0 && len(in.MatchExpressions) == 0 {
		return labels.Nothing(), nil
	}
	return metav1.LabelSelectorAsSelector(&in)
}

// IsEmpty reports whether the selector has no terms.
func IsEmpty(in metav1.LabelSelector) bool {
	return len(in.MatchLabels) == 0 && len(in.MatchExpressions) == 0
}
