# Changelog

All notable changes to HILIOS Operator are documented in this file. The format
is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.5.0] - 2026-05-08

### Added
- Validating webhook stubs for ResiliencePolicy and RecoveryDrill backed by `internal/validation`.
- Drill controller rejects invalid specs at reconcile time with `ValidationFailed` reason.
- Goreleaser ldflags populate the `internal/buildinfo` package; manager logs version on startup.
- Quickstart guide for installing on a kind cluster.
- Drill catalog expanded with Postgres, MySQL, MongoDB, Redis, Elasticsearch, Kafka templates.
- ContentionReport, RebalanceCheck, ResiliencePolicy templates added to the catalog.
- Documentation guide stubs under `docs/guides/` covering 20 common operator topics.
- SonarCloud workflow + project file; quality gate / coverage / maintainability badges in README.
- CodeQL static analysis workflow scheduled weekly.
- PR labeler and stale-issue workflows.

### Changed
- Streamlined Makefile (drop OLM bundle/catalog targets we do not ship).
- Bumped Go to 1.25.10, controller-gen to v0.18.0, golangci-lint to v2.5.0.
- Bumped golang.org/x/net to v0.53.0 to clear GO-2026-4918.

### Fixed
- Data race in healthcheck HTTP runner under parallel tests (mutex guard).
- Dockerfile copies the entire `internal/` tree.

## [0.4.0] - 2026-05-08

### Added
- Cmd and Kubernetes probe types on RecoveryDrill (Litmus-inspired schema).
- HTTP and Kubernetes probes are evaluated in-controller; Pod and Cmd probes
  are recorded as Skip for external runners.
- `internal/healthcheck.RunKubernetes` verifies referenced object existence
  via an unstructured client.
- `internal/safeint` clamps int -> int32 conversions to satisfy gosec G115
  without disabling the rule.
- Drill catalog with Postgres, Kafka, and generic StatefulSet failover templates.
- `docs/probes.md` and `docs/comparison.md` describing the probe schema and
  the project's relationship to Litmus, Velero, and vendor operators.

### Changed
- Streamlined Makefile: removed unused OLM bundle/catalog targets.
- Bumped Go toolchain to 1.25.10 (clears CVE backlog at govulncheck time).
- Bumped controller-gen to v0.18.0 and golangci-lint to v2.5.0.
- Bumped golang.org/x/net to v0.53.0 to clear GO-2026-4918.

### Fixed
- Data race in healthcheck HTTP runner under parallel tests (mutex guard).
- Dockerfile copies the entire `internal/` tree so the builder picks up new
  helper packages.

## [0.3.0] - 2026-05-08

### Added
- ResiliencePolicy joins drill history to compute drift count and report
  freshness violations.
- HTTP health check execution during RecoveryDrill with per-probe evidence.
- New helper packages: `policy/freshness`, `healthcheck`, `mitigation`,
  `contention` (pluggable evaluator).

## [0.2.0] - 2026-05-08

### Added
- Cron-driven re-evaluation via `internal/scheduling`.
- Prometheus collectors for policies, drills, skew, contention findings.
- `hilios.io/paused` annotation honoured at runtime.
- Goreleaser configuration and Helm chart values schema.

### Added
- Project scaffolding via Operator SDK 1.38.
- `resilience.hilios.io/v1alpha1` API group with `ResiliencePolicy`,
  `RecoveryDrill`, `RebalanceCheck`, and `ContentionReport` CRDs.
- Reconcilers for the four CRDs with status conditions, events, and finalizers.
- Unit tests covering the controllers and the internal helper packages.
- Helm chart with deployment, RBAC, CRDs, service, and optional ServiceMonitor.
- CI pipeline (lint, unit, build, helm), Release pipeline, and Security scans.
- Dependabot configuration with auto-merge for minor and patch updates.
- KUTTL end-to-end test suite.

## [0.1.0] - 2026-05-08

Initial preview release. Establishes the public API surface and the developer
workflow. Suitable for evaluation; not yet recommended for production.

### Added
- Project scaffolding via Operator SDK 1.38.
- `resilience.hilios.io/v1alpha1` API group with `ResiliencePolicy`,
  `RecoveryDrill`, `RebalanceCheck`, and `ContentionReport` CRDs.
- Reconcilers for the four CRDs with status conditions, events, and finalizers.
- Unit tests covering the controllers and the internal helper packages.
- Helm chart with deployment, RBAC, CRDs, service, and optional ServiceMonitor.
- CI pipeline (lint, unit, build, helm), Release pipeline, and Security scans.
- Dependabot configuration with auto-merge for minor and patch updates.
- KUTTL end-to-end test suite.
