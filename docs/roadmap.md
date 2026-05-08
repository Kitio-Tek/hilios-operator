# Roadmap

The HILIOS Operator roadmap is anchored on the four pillars described in the
project pitch: **Verification, Protection, Correction, Evidence**. This document
tracks the work in flight and the work planned for upcoming releases.

## Released

### v0.1.0 - Initial Preview
- v1alpha1 API surface with four CRDs.
- Reconcilers with status conditions, finalizers, generation tracking.
- Helm chart with RBAC, leader election, optional ServiceMonitor.
- Unit tests, KUTTL scaffolding, CI pipeline.

### v0.2.0 - Scheduling and Metrics
- Cron-driven re-evaluation via `internal/scheduling`.
- Prometheus collectors for policies, drills, skew, contention findings.
- `hilios.io/paused` annotation honoured at runtime.
- Goreleaser config and Helm chart values schema.

### v0.3.0 - Verification Joins Drill History
- ResiliencePolicy joins drill history to compute drift count.
- HTTP health checks executed during drills with per-probe evidence.
- New helper packages: `policy/freshness`, `healthcheck`, `mitigation`,
  `contention`.

## Planned

### v0.4.0 - Mitigation Engine
- Wire `internal/mitigation` recommendations into RebalanceCheck status.
- Surface authorised mitigations on ResiliencePolicy status.
- Dry-run flag on RecoveryDrill that validates spec without executing.
- NetworkPolicy template in Helm chart (already optional).

### v0.5.0 - Webhook Validation
- Validating admission webhook for ResiliencePolicy and RecoveryDrill.
- Defaulting webhook for cron schedule normalisation.

### v0.6.0 - Velero Integration
- Native Velero `Backup` and `Restore` references on RecoveryDrill.
- Async restore tracking with status condition updates.

### v0.7.0 - Custom Metrics Backends
- `ContentionReport` evaluator backed by Prometheus queries.
- Pluggable `metrics-server` adapter for live throttling signals.

### v0.8.0 - Stable API
- Promote API to `v1` and graduate the operator out of preview.
- Backwards-compatibility guarantees across `v1alpha1` and `v1`.

## Non-goals

- HILIOS does not manage the lifecycle of databases, caches, or message buses.
  Workload-specific operators remain responsible for their own primitives.
- HILIOS does not run user-defined chaos experiments. The drill catalogue is
  scoped to operations a platform team would run on a real outage.
