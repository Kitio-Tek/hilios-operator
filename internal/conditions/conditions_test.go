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

package conditions

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSetAndQuery(t *testing.T) {
	t.Parallel()
	var conds []metav1.Condition
	conds = True(&conds, "Ready", "Reason", "msg", 1)
	if !IsTrue(conds, "Ready") {
		t.Fatal("Ready must be true after True()")
	}
	conds = False(&conds, "Ready", "Reason", "msg", 2)
	if !IsFalse(conds, "Ready") {
		t.Fatal("Ready must be false after False()")
	}
	if IsTrue(conds, "Ready") {
		t.Fatal("Ready must not be reported as true")
	}
}

func TestUnknown(t *testing.T) {
	t.Parallel()
	var conds []metav1.Condition
	conds = Unknown(&conds, "X", "R", "m", 1)
	c := Find(conds, "X")
	if c == nil || c.Status != metav1.ConditionUnknown {
		t.Fatalf("Unknown must set status Unknown, got %#v", c)
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()
	var conds []metav1.Condition
	conds = True(&conds, "X", "R", "m", 1)
	conds = True(&conds, "Y", "R", "m", 1)
	conds = Remove(conds, "X")
	if Find(conds, "X") != nil {
		t.Fatal("Remove must drop the condition")
	}
	if Find(conds, "Y") == nil {
		t.Fatal("Remove must not drop other conditions")
	}
}

func TestSetNilSliceReturnsNil(t *testing.T) {
	t.Parallel()
	if got := Set(nil, "T", "R", "m", metav1.ConditionTrue, 1); got != nil {
		t.Fatalf("Set with nil slice must return nil, got %v", got)
	}
}

func TestSetUpdatesExistingCondition_1(t *testing.T) {
	t.Parallel()
	var conds []metav1.Condition
	conds = True(&conds, "Ready", "First", "msg-1", int64(1))
	conds = False(&conds, "Ready", "Second", "msg-2", int64(1+1))
	c := Find(conds, "Ready")
	if c == nil || c.Status != metav1.ConditionFalse {
		t.Fatalf("update did not transition status: %#v", c)
	}
}

func TestSetUpdatesExistingCondition_2(t *testing.T) {
	t.Parallel()
	var conds []metav1.Condition
	conds = True(&conds, "Ready", "First", "msg-1", int64(2))
	conds = False(&conds, "Ready", "Second", "msg-2", int64(2+1))
	c := Find(conds, "Ready")
	if c == nil || c.Status != metav1.ConditionFalse {
		t.Fatalf("update did not transition status: %#v", c)
	}
}
