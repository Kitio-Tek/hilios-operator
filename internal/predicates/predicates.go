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

// Package predicates provides controller-runtime predicates that suppress
// reconciliation events that the operator does not need to act on. The most
// important predicate is GenerationOrPause, which combines the generation
// changed predicate with a check on the spec.suspend annotation so that
// cluster-wide pauses do not cause reconcile churn.
package predicates

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// GenerationOrPause returns a predicate that triggers reconciliation when the
// generation changes or when the well-known annotation
// "hilios.io/paused" is added or removed.
func GenerationOrPause() predicate.Predicate {
	gen := predicate.GenerationChangedPredicate{}
	return predicate.Or(gen, pauseAnnotationPredicate{})
}

type pauseAnnotationPredicate struct{}

func (pauseAnnotationPredicate) Create(_ event.CreateEvent) bool { return false }
func (pauseAnnotationPredicate) Delete(_ event.DeleteEvent) bool { return false }
func (pauseAnnotationPredicate) Generic(_ event.GenericEvent) bool {
	return false
}

func (pauseAnnotationPredicate) Update(e event.UpdateEvent) bool {
	if e.ObjectOld == nil || e.ObjectNew == nil {
		return false
	}
	const key = "hilios.io/paused"
	oldA := e.ObjectOld.GetAnnotations()[key]
	newA := e.ObjectNew.GetAnnotations()[key]
	return oldA != newA
}
