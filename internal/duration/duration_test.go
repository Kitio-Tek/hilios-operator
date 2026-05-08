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

package duration

import (
	"testing"
	"time"
)

func TestFromSeconds(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   int32
		want time.Duration
	}{
		{0, 0},
		{-1, 0},
		{30, 30 * time.Second},
		{60, time.Minute},
	}
	for _, c := range cases {
		if got := FromSeconds(c.in); got != c.want {
			t.Fatalf("FromSeconds(%d) = %s, want %s", c.in, got, c.want)
		}
	}
}

func TestFromSecondsOr(t *testing.T) {
	t.Parallel()
	if got := FromSecondsOr(0, time.Minute); got != time.Minute {
		t.Fatalf("default not applied: %s", got)
	}
	if got := FromSecondsOr(30, time.Minute); got != 30*time.Second {
		t.Fatalf("value not used: %s", got)
	}
}

func TestClamp(t *testing.T) {
	t.Parallel()
	if got := Clamp(time.Second, 5*time.Second, 10*time.Second); got != 5*time.Second {
		t.Fatalf("low clamp: %s", got)
	}
	if got := Clamp(time.Minute, 5*time.Second, 10*time.Second); got != 10*time.Second {
		t.Fatalf("high clamp: %s", got)
	}
	if got := Clamp(7*time.Second, 5*time.Second, 10*time.Second); got != 7*time.Second {
		t.Fatalf("in-range: %s", got)
	}
}

func TestFromSecondsOrZeroFallsBackToDefault(t *testing.T) {
	t.Parallel()
	d := FromSecondsOr(0, 5*time.Second)
	if d != 5*time.Second {
		t.Fatalf("expected 5s default, got %s", d)
	}
}

func TestClampWithEqualBounds(t *testing.T) {
	t.Parallel()
	if got := Clamp(10*time.Second, 10*time.Second, 10*time.Second); got != 10*time.Second {
		t.Fatalf("equal bounds must pass through, got %s", got)
	}
}

func TestFromSecondsTreatsNegativeAsZero(t *testing.T) {
	t.Parallel()
	if got := FromSeconds(-100); got != 0 {
		t.Fatalf("negative seconds must be zero duration, got %s", got)
	}
}
