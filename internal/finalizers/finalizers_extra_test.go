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

func TestHiliosFinalizerStringValue(t *testing.T) {
	if HiliosFinalizer != "hilios.io/finalizer" {
		t.Fatalf("finalizer string changed: %s", HiliosFinalizer)
	}
}

func TestHasOnNilFinalizers(t *testing.T) {
	t.Parallel()
	obj := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "x"}}
	obj.Finalizers = nil
	if Has(&obj.ObjectMeta, HiliosFinalizer) {
		t.Fatal("nil finalizers must not match")
	}
}
