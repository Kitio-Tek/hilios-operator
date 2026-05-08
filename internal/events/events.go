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

// Package events provides thin wrappers around the controller-runtime event
// recorder. The wrappers exist so that controllers do not have to remember
// the corev1.EventTypeNormal / corev1.EventTypeWarning constants and so that
// the project can later swap the underlying recorder without touching every
// caller.
package events

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// Normal records a Normal-type event on the supplied object.
func Normal(rec record.EventRecorder, obj runtime.Object, reason, messageFmt string, args ...interface{}) {
	if rec == nil || obj == nil {
		return
	}
	rec.Eventf(obj, corev1.EventTypeNormal, reason, messageFmt, args...)
}

// Warning records a Warning-type event on the supplied object.
func Warning(rec record.EventRecorder, obj runtime.Object, reason, messageFmt string, args ...interface{}) {
	if rec == nil || obj == nil {
		return
	}
	rec.Eventf(obj, corev1.EventTypeWarning, reason, messageFmt, args...)
}
