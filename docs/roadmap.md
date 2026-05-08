# Roadmap

The HILIOS Operator roadmap is anchored on the four pillars described in the
project pitch: **Verification, Protection, Correction, Evidence**. This document
tracks the work in flight and the work planned for upcoming releases. Many of
the planned items are inspired by patterns we admire in other operators
(Litmus Chaos, Strimzi, CloudNativePG, Microcks, Vertica).

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
- HTTP health check execution during RecoveryDrill with per-probe evidence.
- Helper packages: `policy/freshness`, `healthcheck`, `mitigation`,
  `contention` (pluggable evaluator).

### v0.4.0 - Probe Schema and Catalog
- Cmd and Kubernetes probe types (Litmus-inspired schema).
- Drill catalog templates for Postgres, MySQL, MongoDB, Redis, Elasticsearch,
  Kafka, and a generic StatefulSet failover template.
- New helper packages: `validation`, `healthcheck/kubernetes`.

### v0.5.0 / v0.5.1 - Project Maturity
- Validating webhook stubs and runtime spec validation.
- Build info exposed via `internal/buildinfo`.
- Documentation guides, ADRs, runbooks, governance docs.
- CodeQL, Trivy, govulncheck, gosec security scans wired into CI.
- Renovate, release-please, stale-issue, labeler workflows.
- Helm chart enhancements: PodDisruptionBudget, NetworkPolicy, values-strict
  and values-ha presets.

## Planned

### v0.6.0 - Velero Integration
- Native Velero `Backup` and `Restore` references on RecoveryDrill.
- Async restore tracking with status condition updates.
- Inspired by [Velero's restore CR pattern](https://velero.io/docs/main/api-types/restore/).

### v0.7.0 - Custom Metrics Backends
- ContentionReport evaluator backed by Prometheus queries (PromQL).
- Pluggable metrics-server adapter for live throttling signals.
- Inspired by [keda](https://keda.sh) external scalers.

### v0.8.0 - Drill Hub
- A read-only ArtifactHub-style index of community drill templates.
- `hilios.io/source` annotation distinguishes hub-installed templates.
- Inspired by [Litmus ChaosHub](https://hub.litmuschaos.io/).

### v0.9.0 - Workflow Orchestration
- Multi-step Drill compositions (chain restore -> failover -> rebalance).
- Inspired by [Argo Workflows](https://github.com/argoproj/argo-workflows) DAG model.

### v0.10.0 - Status Aggregator
- ResilienceReport CRD that summarises one or more policies into a single
  pass/fail dashboard tile.
- Inspired by [Strimzi's KafkaConnect / KafkaConnector relationship](https://github.com/strimzi/strimzi-kafka-operator).

### v0.11.0 - Backup Engine Adapter Set
- Adapters for Restic, Kasten K10, Stash, in addition to Velero.
- Inspired by the [CloudNativePG plugin model](https://cloudnative-pg.io/documentation/current/architecture/).

### v0.12.0 - Mock Workload Harness
- A small fixture engine that creates synthetic workloads for local kind
  testing.
- Inspired by [Microcks](https://microcks.io/) mock workflow generation.

### v1.0.0 - Stable API
- Promote API to `v1` and graduate the operator out of preview.
- Backwards-compatibility guarantees across `v1alpha1` and `v1`.
- Conformance test suite.

## Non-goals

- HILIOS does not manage the lifecycle of databases, caches, or message buses.
  Workload-specific operators remain responsible for their own primitives.
- HILIOS does not run user-defined chaos experiments. The drill catalogue is
  scoped to operations a platform team would run on a real outage.
- HILIOS does not build a bespoke metrics pipeline. It plugs into Prometheus
  and the metrics-server APIs that already exist in the cluster.

## How to influence the roadmap

- File a GitHub issue with the `enhancement` label and a clear motivation.
- Propose an ADR (`docs/adr/`) for a structural change.
- Open a draft PR with a working prototype - the smaller the first PR the
  faster it lands.
