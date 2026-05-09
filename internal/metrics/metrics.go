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

// Package metrics declares the Prometheus collectors emitted by the HILIOS
// controllers. The collectors are registered against the controller-runtime
// metrics registry through MustRegister so that they appear under the
// controller-manager metrics endpoint.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

const subsystem = "hilios"

var (
	// PolicyEvaluationsTotal counts how many times each ResiliencePolicy has
	// been evaluated since the manager started.
	PolicyEvaluationsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: subsystem,
		Name:      "policy_evaluations_total",
		Help:      "Total number of ResiliencePolicy evaluations.",
	}, []string{"namespace", "name"})

	// PolicyDriftCount records the number of violations observed in the most
	// recent evaluation of each ResiliencePolicy.
	PolicyDriftCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: subsystem,
		Name:      "policy_drift_count",
		Help:      "Number of policy violations observed in the most recent evaluation.",
	}, []string{"namespace", "name"})

	// DrillDurationSeconds records the wall-clock duration of completed drills.
	DrillDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: subsystem,
		Name:      "drill_duration_seconds",
		Help:      "Wall-clock duration of completed RecoveryDrills.",
		Buckets:   prometheus.ExponentialBuckets(1, 2, 12),
	}, []string{"type", "result"})

	// DrillsTotal counts completed drills by type and outcome.
	DrillsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: subsystem,
		Name:      "drills_total",
		Help:      "Total number of RecoveryDrills that reached a terminal phase.",
	}, []string{"type", "result"})

	// RebalanceSkew records the most recent skew observed by a RebalanceCheck.
	RebalanceSkew = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: subsystem,
		Name:      "rebalance_skew",
		Help:      "Most recent topology skew observed by a RebalanceCheck.",
	}, []string{"namespace", "name"})

	// ContentionFindings records the number of contention findings produced
	// by the most recent evaluation of a ContentionReport.
	ContentionFindings = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: subsystem,
		Name:      "contention_findings",
		Help:      "Number of contention findings produced by the most recent evaluation.",
	}, []string{"namespace", "name"})
)

// init registers the collectors with the controller-runtime metrics registry.
func init() {
	metrics.Registry.MustRegister(
		PolicyEvaluationsTotal,
		PolicyDriftCount,
		DrillDurationSeconds,
		DrillsTotal,
		RebalanceSkew,
		ContentionFindings,
	)
}
