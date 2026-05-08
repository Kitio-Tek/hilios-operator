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

// Package healthcheck implements the small probe runners used by RecoveryDrill
// when verifying that a restored or failed-over workload is healthy. The
// HTTP runner is exposed as a package-level variable so callers can replace it
// in tests.
package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
)

// HTTPRunner runs the HTTP-flavoured health check against a URL and returns
// the observed status code. Replace via SetHTTPRunner in tests.
type HTTPRunner func(ctx context.Context, url string, timeout time.Duration) (int, error)

var (
	httpRunnerMu sync.RWMutex
	httpRunner   HTTPRunner = defaultHTTPRunner
)

// SetHTTPRunner installs a custom HTTP runner. Pass nil to restore the default.
// The variable is mutex-protected so concurrent tests can swap it without
// triggering the race detector.
func SetHTTPRunner(r HTTPRunner) {
	httpRunnerMu.Lock()
	defer httpRunnerMu.Unlock()
	if r == nil {
		httpRunner = defaultHTTPRunner
		return
	}
	httpRunner = r
}

func getHTTPRunner() HTTPRunner {
	httpRunnerMu.RLock()
	defer httpRunnerMu.RUnlock()
	return httpRunner
}

// Run executes the health check and returns nil on success.
func Run(ctx context.Context, hc resiliencev1alpha1.HealthCheck) error {
	timeout := time.Duration(hc.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	switch hc.Type {
	case "HTTP":
		if hc.URL == "" {
			return fmt.Errorf("http health check %q has empty URL", hc.Name)
		}
		expected := hc.ExpectedStatusCode
		if expected == 0 {
			expected = http.StatusOK
		}
		code, err := getHTTPRunner()(ctx, hc.URL, timeout)
		if err != nil {
			return fmt.Errorf("http health check %q: %w", hc.Name, err)
		}
		if int(expected) != code {
			return fmt.Errorf("http health check %q: status %d, expected %d", hc.Name, code, expected)
		}
		return nil
	case "Pod":
		// Pod probes are evaluated by the controller against the API server.
		// They are short-circuited here so that callers can pre-validate the
		// spec without performing a live read.
		return nil
	default:
		return fmt.Errorf("unknown health check type %q", hc.Type)
	}
}

func defaultHTTPRunner(ctx context.Context, url string, timeout time.Duration) (int, error) {
	c := &http.Client{Timeout: timeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()
	return resp.StatusCode, nil
}
