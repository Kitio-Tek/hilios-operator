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

package cronexpr

import (
	"testing"
	"time"
)

func TestParseEmptyOrInvalid(t *testing.T) {
	t.Parallel()
	cases := []string{"", "not-a-cron", "99 * * * *", "0 25 * * *", "0 0 32 * *", "0 0 * * 9"}
	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			if _, err := Parse(expr); err == nil {
				t.Fatalf("Parse(%q) must return an error", expr)
			}
		})
	}
}

func TestParseValid(t *testing.T) {
	t.Parallel()
	cases := []string{
		"* * * * *",
		"*/5 * * * *",
		"0 * * * *",
		"0 0 * * *",
		"0 0 1 * *",
		"0 9 * * 1-5",
		"0 0 * * 0",
	}
	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			if _, err := Parse(expr); err != nil {
				t.Fatalf("Parse(%q) failed: %v", expr, err)
			}
		})
	}
}

func TestNextEveryFiveMinutes(t *testing.T) {
	t.Parallel()
	s, err := Parse("*/5 * * * *")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	from := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	next := s.Next(from)
	want := time.Date(2026, time.January, 1, 0, 5, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Fatalf("Next: got %s want %s", next, want)
	}
}

func TestNextHourlyRollsForward(t *testing.T) {
	t.Parallel()
	s, _ := Parse("0 * * * *")
	from := time.Date(2026, time.January, 1, 5, 30, 0, 0, time.UTC)
	if got := s.Next(from).Hour(); got != 6 {
		t.Fatalf("Next 0 * * * * from 5:30 should be 6, got %d", got)
	}
}

func TestNextWeeklyHitsSunday(t *testing.T) {
	t.Parallel()
	s, _ := Parse("0 0 * * 0")
	from := time.Date(2026, time.January, 5, 0, 0, 0, 0, time.UTC)
	if got := s.Next(from).Weekday(); got != time.Sunday {
		t.Fatalf("Next should be Sunday, got %s", got)
	}
}

func TestExprRoundTrip(t *testing.T) {
	t.Parallel()
	expr := "0 * * * *"
	s, _ := Parse(expr)
	if s.Expr() != expr {
		t.Fatalf("Expr: got %s want %s", s.Expr(), expr)
	}
}
