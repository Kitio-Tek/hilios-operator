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
		in   int
		want int32
	}{
		{0, 0},
		{-5, 0},
		{42, 42},
		{math.MaxInt32, math.MaxInt32},
		{math.MaxInt32 + 1, math.MaxInt32},
	}
	for _, c := range cases {
		if got := Int32(c.in); got != c.want {
			t.Fatalf("Int32(%d) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestInt32From64(t *testing.T) {
	t.Parallel()
	if got := Int32From64(int64(math.MaxInt32) + 100); got != math.MaxInt32 {
		t.Fatalf("clamp expected, got %d", got)
	}
	if got := Int32From64(-1); got != 0 {
		t.Fatalf("negative expected 0, got %d", got)
	}
	if got := Int32From64(7); got != 7 {
		t.Fatalf("passthrough expected 7, got %d", got)
	}
}

func TestInt32EqualsInt(t *testing.T) {
	t.Parallel()
	if !Int32EqualsInt(int32(5), 5) {
		t.Fatal("equal must compare true")
	}
	if Int32EqualsInt(int32(5), 6) {
		t.Fatal("unequal must compare false")
	}
}

func TestInt32MaxBoundary(t *testing.T) {
	t.Parallel()
	if got := Int32(math.MaxInt32); got != math.MaxInt32 {
		t.Fatalf("MaxInt32 must pass through, got %d", got)
	}
}

func TestInt32From64NegativeFloor(t *testing.T) {
	t.Parallel()
	if got := Int32From64(-100); got != 0 {
		t.Fatalf("negative input must clamp to zero, got %d", got)
	}
}
