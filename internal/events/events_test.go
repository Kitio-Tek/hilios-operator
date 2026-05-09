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

package events

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
)

func TestNormalAndWarningNoOpOnNilRecorder(t *testing.T) {
	t.Parallel()
	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "p"}}
	Normal(nil, pod, "X", "msg")
	Warning(nil, pod, "X", "msg")
}

func TestNormalAndWarningEmitOnFakeRecorder(t *testing.T) {
	t.Parallel()
	rec := record.NewFakeRecorder(4)
	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "p"}}
	Normal(rec, pod, "Reason", "msg %s", "x")
	Warning(rec, pod, "Reason", "msg %s", "y")
	for i := 0; i < 2; i++ {
		select {
		case ev := <-rec.Events:
			if ev == "" {
				t.Fatalf("expected a non-empty event, got %q", ev)
			}
		default:
			t.Fatalf("expected at least %d events on the channel", i+1)
		}
	}
}

func TestNormalNoOpOnNilObject(t *testing.T) {
	t.Parallel()
	rec := record.NewFakeRecorder(1)
	Normal(rec, nil, "Reason", "msg")
	select {
	case ev := <-rec.Events:
		t.Fatalf("expected no event for nil object, got %q", ev)
	default:
	}
}
