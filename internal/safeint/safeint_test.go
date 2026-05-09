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

package safeint

import (
	"math"
	"testing"
)

func TestInt32(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   int
		want int32
	}{
		{"zero", 0, 0},
		{"negative", -5, 0},
		{"in-range", 42, 42},
		{"max-boundary", math.MaxInt32, math.MaxInt32},
		{"overflow", math.MaxInt32 + 1, math.MaxInt32},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Int32(c.in); got != c.want {
				t.Fatalf("Int32(%d) = %d, want %d", c.in, got, c.want)
			}
		})
	}
}

func TestInt32From64(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   int64
		want int32
	}{
		{"overflow", int64(math.MaxInt32) + 100, math.MaxInt32},
		{"negative", -1, 0},
		{"in-range", 7, 7},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Int32From64(c.in); got != c.want {
				t.Fatalf("Int32From64(%d) = %d, want %d", c.in, got, c.want)
			}
		})
	}
}
