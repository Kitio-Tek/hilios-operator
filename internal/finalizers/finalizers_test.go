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

package finalizers

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestHasFalseWhenAbsent(t *testing.T) {
	t.Parallel()
	obj := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "x"}}
	if Has(&obj.ObjectMeta, HiliosFinalizer) {
		t.Fatal("Has must return false on empty finalizer list")
	}
}

func TestAddAndRemove(t *testing.T) {
	t.Parallel()
	obj := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "x"}}

	if !Add(&obj.ObjectMeta, HiliosFinalizer) {
		t.Fatal("Add must return true on first add")
	}
	if Add(&obj.ObjectMeta, HiliosFinalizer) {
		t.Fatal("Add must return false when already present")
	}
	if !Has(&obj.ObjectMeta, HiliosFinalizer) {
		t.Fatal("Has must return true after Add")
	}
	if !Remove(&obj.ObjectMeta, HiliosFinalizer) {
		t.Fatal("Remove must return true on existing")
	}
	if Remove(&obj.ObjectMeta, HiliosFinalizer) {
		t.Fatal("Remove must return false when absent")
	}
}
