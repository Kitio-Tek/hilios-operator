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

package labels

import "testing"

func TestManagedBy(t *testing.T) {
	t.Parallel()
	got := ManagedBy()
	if got[LabelManagedBy] != LabelManagedByValue {
		t.Fatalf("managed-by label missing, got %#v", got)
	}
}

func TestMergeManagedByPreservesInputs(t *testing.T) {
	t.Parallel()
	in := map[string]string{"app": "x", LabelManagedBy: "stale"}
	got := MergeManagedBy(in)
	if got["app"] != "x" {
		t.Fatal("input keys must be preserved")
	}
	if got[LabelManagedBy] != LabelManagedByValue {
		t.Fatal("managed-by must be overwritten")
	}
	// Input map must not be mutated.
	if in[LabelManagedBy] != "stale" {
		t.Fatal("MergeManagedBy must not mutate input")
	}
}

func TestMergeManagedByEmptyInput(t *testing.T) {
	t.Parallel()
	got := MergeManagedBy(nil)
	if got[LabelManagedBy] != LabelManagedByValue {
		t.Fatal("nil input must yield managed-by entry")
	}
}

func TestLabelConstantsHaveHiliosPrefix(t *testing.T) {
	t.Parallel()
	cases := []string{LabelEnabled, LabelPolicy, LabelDrill, LabelCheck, LabelReport}
	for _, c := range cases {
		if c == "" || c[:7] != "hilios." {
			t.Fatalf("label %q must start with hilios.", c)
		}
	}
}

func TestSetWithOwner(t *testing.T) {
	t.Parallel()
	got := Set(LabelPolicy, "demo")
	if got[LabelPolicy] != "demo" || got[LabelManagedBy] != LabelManagedByValue {
		t.Fatalf("unexpected map: %#v", got)
	}
}
