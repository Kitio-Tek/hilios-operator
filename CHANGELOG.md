# Changelog

All notable changes to HILIOS Operator are documented in this file. The format
is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
