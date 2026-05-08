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

package healthcheck

import (
	"context"
	"errors"
	"testing"
	"time"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
)

func TestRunPodAlwaysOK(t *testing.T) {
	t.Parallel()
	if err := Run(context.Background(), resiliencev1alpha1.HealthCheck{Name: "p", Type: "Pod"}); err != nil {
		t.Fatalf("Pod type must short-circuit ok, got %v", err)
	}
}

func TestRunHTTPSuccess(t *testing.T) {
	t.Parallel()
	SetHTTPRunner(func(_ context.Context, url string, _ time.Duration) (int, error) {
		if url != "http://example/healthz" {
			t.Fatalf("unexpected url: %s", url)
		}
		return 200, nil
	})
	t.Cleanup(func() { SetHTTPRunner(nil) })

	hc := resiliencev1alpha1.HealthCheck{Name: "h", Type: "HTTP", URL: "http://example/healthz"}
	if err := Run(context.Background(), hc); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestRunHTTPMismatchedStatus(t *testing.T) {
	t.Parallel()
	SetHTTPRunner(func(_ context.Context, _ string, _ time.Duration) (int, error) {
		return 500, nil
	})
	t.Cleanup(func() { SetHTTPRunner(nil) })

	hc := resiliencev1alpha1.HealthCheck{Name: "h", Type: "HTTP", URL: "http://example/healthz"}
	if err := Run(context.Background(), hc); err == nil {
		t.Fatal("expected error on 500")
	}
}

func TestRunHTTPRunnerError(t *testing.T) {
	t.Parallel()
	SetHTTPRunner(func(_ context.Context, _ string, _ time.Duration) (int, error) {
		return 0, errors.New("transport failed")
	})
	t.Cleanup(func() { SetHTTPRunner(nil) })

	hc := resiliencev1alpha1.HealthCheck{Name: "h", Type: "HTTP", URL: "http://example/healthz"}
	if err := Run(context.Background(), hc); err == nil {
		t.Fatal("expected error from runner")
	}
}

func TestRunUnknownType(t *testing.T) {
	t.Parallel()
	if err := Run(context.Background(), resiliencev1alpha1.HealthCheck{Name: "x", Type: "TCP"}); err == nil {
		t.Fatal("expected error for unknown type")
	}
}
