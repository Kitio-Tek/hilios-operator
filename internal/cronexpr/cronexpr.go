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

// Package cronexpr is a thin wrapper around robfig/cron used by the
// scheduling logic. The wrapper is in place so that controllers depend on a
// stable interface and the underlying cron implementation can be swapped
// without touching every reconciler.
package cronexpr

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// Schedule represents a parsed cron expression and its lazily computed
// activation times. Instances are immutable once Parse returns.
type Schedule struct {
	expr string
	sch  cron.Schedule
}

// Parse parses a five-field cron expression in the standard cron format.
// Empty input returns an error.
func Parse(expr string) (*Schedule, error) {
	if expr == "" {
		return nil, fmt.Errorf("empty cron expression")
	}
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	s, err := parser.Parse(expr)
	if err != nil {
		return nil, err
	}
	return &Schedule{expr: expr, sch: s}, nil
}

// Next returns the next activation time strictly after from.
func (s *Schedule) Next(from time.Time) time.Time {
	return s.sch.Next(from)
}

// Expr returns the original cron string.
func (s *Schedule) Expr() string {
	return s.expr
}
