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

// Package scheduling computes the next re-evaluation time for resources that
// declare an optional cron schedule. It is intentionally pure so reconcilers
// can test scheduling decisions without instantiating a real clock.
package scheduling

import (
	"time"

	"github.com/Kitio-Tek/hilios-operator/internal/cronexpr"
)

// NextRequeue returns the duration the controller should requeue after.
// Empty or invalid expressions fall back to the supplied default duration so
// callers do not have to pre-validate their cron strings.
// expr is the cron expression (may be empty), now is the controller clock,
// and fallback is the duration used when expr is empty or invalid.
func NextRequeue(expr string, now time.Time, fallback time.Duration) time.Duration {
	if expr == "" {
		return fallback
	}
	s, err := cronexpr.Parse(expr)
	if err != nil {
		return fallback
	}
	d := s.Next(now).Sub(now)
	if d <= 0 {
		return fallback
	}
	return d
}
