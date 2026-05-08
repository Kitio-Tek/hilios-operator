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

func TestParseEmpty(t *testing.T) {
	t.Parallel()
	if _, err := Parse(""); err == nil {
		t.Fatal("expected error for empty cron")
	}
}

func TestParseInvalid(t *testing.T) {
	t.Parallel()
	if _, err := Parse("not-a-cron"); err == nil {
		t.Fatal("expected error for invalid cron")
	}
}

func TestParseAndNext(t *testing.T) {
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

func TestExpr(t *testing.T) {
	t.Parallel()
	s, _ := Parse("0 * * * *")
	if s.Expr() != "0 * * * *" {
		t.Fatalf("Expr: got %s", s.Expr())
	}
}

func TestParseWeeklySchedule(t *testing.T) {
	t.Parallel()
	s, err := Parse("0 0 * * 0")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	from := time.Date(2026, time.January, 5, 0, 0, 0, 0, time.UTC)
	next := s.Next(from)
	if next.Weekday() != time.Sunday {
		t.Fatalf("Next should be Sunday, got %s", next.Weekday())
	}
}

func TestExprEmpty(t *testing.T) {
	t.Parallel()
	if _, err := Parse(""); err == nil {
		t.Fatal("Parse must reject empty expression")
	}
}

func TestParseHourly(t *testing.T) {
	t.Parallel()
	s, err := Parse("0 * * * *")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	from := time.Date(2026, time.January, 1, 5, 30, 0, 0, time.UTC)
	if got := s.Next(from); got.Hour() != 6 {
		t.Fatalf("Next 0 * * * * from 5:30 should be 6:00, got %s", got)
	}
}

func TestValidateOK(t *testing.T) {
	t.Parallel()
	if err := Validate("*/5 * * * *"); err != nil {
		t.Fatalf("Validate must accept valid expr: %v", err)
	}
}

func TestValidateError(t *testing.T) {
	t.Parallel()
	if err := Validate(""); err == nil {
		t.Fatal("Validate must reject empty expr")
	}
}

func TestParseAndExpr_1(t *testing.T) {
	t.Parallel()
	expr := "1 * * * *"
	s, err := Parse(expr)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if s.Expr() != expr {
		t.Fatalf("Expr roundtrip: got %s want %s", s.Expr(), expr)
	}
}

func TestParseAndExpr_2(t *testing.T) {
	t.Parallel()
	expr := "2 * * * *"
	s, err := Parse(expr)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if s.Expr() != expr {
		t.Fatalf("Expr roundtrip: got %s want %s", s.Expr(), expr)
	}
}

func TestParseAndExpr_3(t *testing.T) {
	t.Parallel()
	expr := "3 * * * *"
	s, err := Parse(expr)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if s.Expr() != expr {
		t.Fatalf("Expr roundtrip: got %s want %s", s.Expr(), expr)
	}
}

func TestParseAndExpr_4(t *testing.T) {
	t.Parallel()
	expr := "4 * * * *"
	s, err := Parse(expr)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if s.Expr() != expr {
		t.Fatalf("Expr roundtrip: got %s want %s", s.Expr(), expr)
	}
}
