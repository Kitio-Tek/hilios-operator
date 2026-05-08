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

package predicates

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func obj(gen int64, ann map[string]string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Generation:  gen,
			Annotations: ann,
		},
	}
}

func TestGenerationOrPauseUpdate(t *testing.T) {
	t.Parallel()
	p := GenerationOrPause()

	cases := []struct {
		name   string
		old    *corev1.ConfigMap
		newObj *corev1.ConfigMap
		want   bool
	}{
		{"no change", obj(1, nil), obj(1, nil), false},
		{"generation changed", obj(1, nil), obj(2, nil), true},
		{"pause added", obj(1, nil), obj(1, map[string]string{"hilios.io/paused": "true"}), true},
		{"pause unchanged", obj(1, map[string]string{"hilios.io/paused": "true"}), obj(1, map[string]string{"hilios.io/paused": "true"}), false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := p.Update(event.UpdateEvent{ObjectOld: c.old, ObjectNew: c.newObj})
			if got != c.want {
				t.Fatalf("Update() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestGenerationOrPauseDelegatesCreate(t *testing.T) {
	t.Parallel()
	p := GenerationOrPause()

	// GenerationChangedPredicate returns true for Create events (controller-runtime default),
	// so Or(GenerationChanged, PauseAnnotation) must too.
	if !p.Create(event.CreateEvent{Object: obj(1, nil)}) {
		t.Fatalf("expected Create to be true")
	}
}
