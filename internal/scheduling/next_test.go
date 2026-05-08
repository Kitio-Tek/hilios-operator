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

package scheduling

import (
	"testing"
	"time"
)

func TestNextRequeueEmpty(t *testing.T) {
	t.Parallel()
	got := NextRequeue("", time.Now(), 5*time.Minute)
	if got != 5*time.Minute {
		t.Fatalf("empty schedule fallback want 5m, got %s", got)
	}
}

func TestNextRequeueInvalid(t *testing.T) {
	t.Parallel()
	got := NextRequeue("not-a-cron", time.Now(), 5*time.Minute)
	if got != 5*time.Minute {
		t.Fatalf("invalid schedule fallback want 5m, got %s", got)
	}
}

func TestNextRequeueValid(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC)
	got := NextRequeue("*/15 * * * *", now, time.Minute)
	want := 15 * time.Minute
	if got != want {
		t.Fatalf("valid schedule want %s, got %s", want, got)
	}
}

func TestNextRequeueLessThanMinute(t *testing.T) {
	t.Parallel()
	got := NextRequeue("* * * * *", time.Now(), 5*time.Minute)
	if got > time.Minute {
		t.Fatalf("every-minute schedule must yield <= 1m, got %s", got)
	}
}

func TestNextRequeueEveryNMinutes_1(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC)
	got := NextRequeue("*/1 * * * *", now, time.Second)
	want := time.Duration(1) * time.Minute
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestNextRequeueEveryNMinutes_5(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC)
	got := NextRequeue("*/5 * * * *", now, time.Second)
	want := time.Duration(5) * time.Minute
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestNextRequeueEveryNMinutes_10(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC)
	got := NextRequeue("*/10 * * * *", now, time.Second)
	want := time.Duration(10) * time.Minute
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestNextRequeueEveryNMinutes_15(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC)
	got := NextRequeue("*/15 * * * *", now, time.Second)
	want := time.Duration(15) * time.Minute
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}
