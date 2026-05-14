# Hilios Operator

[![CI](https://github.com/Kitio-Tek/hilios-operator/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Kitio-Tek/hilios-operator/actions/workflows/ci.yml)
[![E2E](https://github.com/Kitio-Tek/hilios-operator/actions/workflows/e2e.yml/badge.svg?branch=main)](https://github.com/Kitio-Tek/hilios-operator/actions/workflows/e2e.yml)
[![Security](https://github.com/Kitio-Tek/hilios-operator/actions/workflows/security.yml/badge.svg?branch=main)](https://github.com/Kitio-Tek/hilios-operator/actions/workflows/security.yml)
[![CodeQL](https://github.com/Kitio-Tek/hilios-operator/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/Kitio-Tek/hilios-operator/actions/workflows/codeql.yml)
[![Trivy](https://github.com/Kitio-Tek/hilios-operator/actions/workflows/trivy.yml/badge.svg?branch=main)](https://github.com/Kitio-Tek/hilios-operator/actions/workflows/trivy.yml)
[![Latest release](https://img.shields.io/github/v/release/Kitio-Tek/hilios-operator?sort=semver)](https://github.com/Kitio-Tek/hilios-operator/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/Kitio-Tek/hilios-operator.svg)](https://pkg.go.dev/github.com/Kitio-Tek/hilios-operator)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kitio-Tek/hilios-operator)](https://goreportcard.com/report/github.com/Kitio-Tek/hilios-operator)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](LICENSE)

Hilios Operator is a Kubernetes operator for resilience enforcement and corrective
orchestration of distributed workloads. It continuously evaluates declarative
resilience policies and executes guarded corrective workflows such as restore
verification, topology rebalancing checks, and noisy-neighbor mitigation while
recording audit-friendly evidence as Kubernetes-native status conditions.

## Why Hilios

Backups, failover plans, and workload fairness controls often exist only on paper.
Hilios treats resilience as a control loop: a `ResiliencePolicy` declares what
should be true about a workload, the controller continuously evaluates the
declaration, and `RecoveryDrill`, `RebalanceCheck`, and `ContentionReport`
resources record what was tested and what was found. Operators do not deploy a
new database, message bus, or pipeline; they enforce the operational invariants
that protect the workloads already running in the cluster.

## Installation

### Prerequisites

- Kubernetes 1.27 or later
- Helm 3.x
- kubectl

### Installing with Helm

The Helm chart is published to a GitHub Pages chart repository on every
push to `main` that touches `charts/`. Install the latest release with:

```bash
helm repo add hilios https://kitio-tek.github.io/hilios-operator
helm repo update
helm install hilios hilios/hilios-operator \
  --namespace hilios-system \
  --create-namespace
```

To install directly from a working tree (development mode):

```bash
helm install hilios charts/hilios-operator/ \
  --namespace hilios-system \
  --create-namespace
```

Verify the operator pod is running:

```bash
kubectl get pods -n hilios-system
```

The chart installs the four CRDs (`ResiliencePolicy`, `RecoveryDrill`,
`RebalanceCheck`, `ContentionReport`), the manager Deployment, the cluster role
binding, and an optional `ServiceMonitor` for Prometheus.

Once installed, see [docs/usage.md](docs/usage.md) for declaring a
`ResiliencePolicy` and running drills, checks, and reports.

## Features

- Declarative ResiliencePolicy with target selector, verifications, and mitigations
- RecoveryDrill state machine with audit-friendly evidence trail
- RebalanceCheck across arbitrary topology keys
- ContentionReport with pluggable signal evaluator (default: PodConditions)
- Standard Kubernetes API: conditions, events, finalizers, generation tracking
- Helm chart with RBAC, leader election, ServiceMonitor
- [KUTTL](https://github.com/kudobuilder/kuttl) and [Chainsaw](https://github.com/kyverno/chainsaw) end-to-end tests
- Dependabot with auto-merge for minor/patch updates
- govulncheck, gosec, and Trivy filesystem scans on every PR
- CodeQL weekly Go static analysis

## Documentation

- [Usage](docs/usage.md) - declaring policies and running drills, checks, and reports.
- [API reference](docs/api-reference.md) - field tables for the four CRDs.
- [Architecture](docs/architecture.md) - reconciler topology and internal packages.
- [Operations](docs/operations.md) - logging, ports, RBAC, leader election, telemetry, and the kubectl helper.
- [Conditions reference](docs/conditions.md) - every condition the operator emits.
- [Metrics reference](docs/metrics.md) - Prometheus collectors and how to scrape them.
- [Probes](docs/probes.md) - HTTP, Pod, Cmd, and Kubernetes probe semantics.
- [Comparison](docs/comparison.md) - relationship to Litmus, Velero, vendor operators.
- [Roadmap](docs/roadmap.md) - planned releases.
- [Drill catalog](examples/catalog/README.md) - reusable RecoveryDrill templates.

## Development

See [DEVELOPER.md](DEVELOPER.md) for local development setup, including how to
run the operator against a kind cluster and how to run unit, integration,
KUTTL, and Chainsaw tests.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

## Contributors

Thanks to everyone who has filed issues, opened pull requests, or pushed
code to this project. Run `git shortlog -sn --no-merges` for a current list.

## License

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
