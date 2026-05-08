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

package topology

import "testing"

func TestSkewWith2by3Imbalance(t *testing.T) {
	d := Distribution{"a": 2, "b": 5}
	if got := Skew(d); got != 3 {
		t.Fatalf("expected skew 3, got %d", got)
	}
}

func TestSkewAcrossManyZones(t *testing.T) {
	d := Distribution{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}
	if got := Skew(d); got != 4 {
		t.Fatalf("expected skew 4, got %d", got)
	}
}
