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

// Package finalizers provides constants and helpers for managing the
// finalizer strings used across HILIOS controllers.
package finalizers

import "k8s.io/apimachinery/pkg/api/meta"

// HiliosFinalizer is the canonical finalizer applied to HILIOS resources that
// require asynchronous cleanup before deletion completes.
const HiliosFinalizer = "hilios.io/finalizer"

// Has reports whether the supplied object metadata has the finalizer.
func Has(obj metav1Object, name string) bool {
	for _, f := range obj.GetFinalizers() {
		if f == name {
			return true
		}
	}
	return false
}

// Add appends name to the object's finalizer list if not already present and
// returns true when the slice was modified.
func Add(obj metav1Object, name string) bool {
	for _, f := range obj.GetFinalizers() {
		if f == name {
			return false
		}
	}
	obj.SetFinalizers(append(obj.GetFinalizers(), name))
	return true
}

// Remove drops name from the object's finalizer list and returns true when the
// slice was modified.
func Remove(obj metav1Object, name string) bool {
	in := obj.GetFinalizers()
	out := make([]string, 0, len(in))
	removed := false
	for _, f := range in {
		if f == name {
			removed = true
			continue
		}
		out = append(out, f)
	}
	if removed {
		obj.SetFinalizers(out)
	}
	return removed
}

// metav1Object is the minimal subset of metav1.Object the helpers use.
type metav1Object interface {
	GetFinalizers() []string
	SetFinalizers([]string)
}

// AccessorFor delegates to meta.Accessor and returns just the GetFinalizers /
// SetFinalizers slice for the supplied object. It is exposed so that callers
// who hold a runtime.Object (rather than the concrete metav1.Object) can use
// the same helpers without importing meta directly.
func AccessorFor(obj interface{}) ([]string, error) {
	a, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}
	return a.GetFinalizers(), nil
}
