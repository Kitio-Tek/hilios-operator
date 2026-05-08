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

package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestPolicyEvaluationsCounter(t *testing.T) {
	t.Parallel()
	PolicyEvaluationsTotal.Reset()
	PolicyEvaluationsTotal.WithLabelValues("ns", "p").Inc()
	PolicyEvaluationsTotal.WithLabelValues("ns", "p").Inc()
	if got := testutil.ToFloat64(PolicyEvaluationsTotal.WithLabelValues("ns", "p")); got != 2 {
		t.Fatalf("counter want 2, got %v", got)
	}
}

func TestRebalanceSkewGauge(t *testing.T) {
	t.Parallel()
	RebalanceSkew.Reset()
	RebalanceSkew.WithLabelValues("ns", "c").Set(3)
	if got := testutil.ToFloat64(RebalanceSkew.WithLabelValues("ns", "c")); got != 3 {
		t.Fatalf("gauge want 3, got %v", got)
	}
}

func TestPolicyDriftCounterReset(t *testing.T) {
	t.Parallel()
	PolicyDriftCount.Reset()
	PolicyDriftCount.WithLabelValues("ns", "p").Set(7)
	if got := testutil.ToFloat64(PolicyDriftCount.WithLabelValues("ns", "p")); got != 7 {
		t.Fatalf("gauge want 7, got %v", got)
	}
}

func TestContentionFindingsGauge(t *testing.T) {
	t.Parallel()
	ContentionFindings.Reset()
	ContentionFindings.WithLabelValues("ns", "r").Set(3)
	if got := testutil.ToFloat64(ContentionFindings.WithLabelValues("ns", "r")); got != 3 {
		t.Fatalf("gauge want 3, got %v", got)
	}
}
