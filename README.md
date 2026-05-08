# HILIOS Operator

[![CI](https://github.com/Kitio-Tek/hilios-operator/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Kitio-Tek/hilios-operator/actions/workflows/ci.yml)
[![Latest release](https://img.shields.io/github/v/release/Kitio-Tek/hilios-operator?sort=semver)](https://github.com/Kitio-Tek/hilios-operator/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kitio-Tek/hilios-operator)](https://goreportcard.com/report/github.com/Kitio-Tek/hilios-operator)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](LICENSE)

HILIOS Operator is a Kubernetes operator for resilience enforcement and corrective
orchestration of distributed workloads. It continuously evaluates declarative
resilience policies and executes guarded corrective workflows such as restore
verification, topology rebalancing checks, and noisy-neighbor mitigation while
recording audit-friendly evidence as Kubernetes-native status conditions.

## Why HILIOS

Backups, failover plans, and workload fairness controls often exist only on paper.
HILIOS treats resilience as a control loop: a `ResiliencePolicy` declares what
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

Install the operator into a dedicated namespace:

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

## Configuration

### Declaring a ResiliencePolicy

A `ResiliencePolicy` selects workloads by label and authorises the verifications
and mitigations the controller may run against them. Workloads opt into
governance by carrying the label `hilios.io/enabled=true`.

```yaml
apiVersion: resilience.hilios.io/v1alpha1
kind: ResiliencePolicy
metadata:
  name: payments
spec:
  targetSelector:
    matchLabels:
      hilios.io/enabled: "true"
      tier: critical
  verifications:
    - kind: RestoreVerification
      intervalSeconds: 21600
      freshnessSeconds: 86400
    - kind: ReplicaPlacement
      intervalSeconds: 600
      freshnessSeconds: 3600
  mitigations:
    - ApplyTopologySpread
    - ScaleSafely
  slo:
    recoveryTimeSeconds: 600
    maxReplicaSkew: 1
```

```bash
kubectl apply -f policy.yaml
kubectl get rpol
kubectl describe rpol payments
```

### Running a RecoveryDrill

A `RecoveryDrill` is a one-shot or scheduled verification. The current build
covers two drill types: `RestoreVerification` (used to validate that a backup
can be restored into a temporary namespace) and `FailoverDrill` (used to
exercise the failover path of a stateful workload).

```yaml
apiVersion: resilience.hilios.io/v1alpha1
kind: RecoveryDrill
metadata:
  name: payments-restore
spec:
  type: RestoreVerification
  source:
    name: payments-2026-05-07
    namespace: velero
  verificationNamespace: payments-restore-verify
  healthChecks:
    - name: pod-readiness
      type: Pod
      podSelector:
        matchLabels:
          app: payments
      timeoutSeconds: 60
  cleanup: true
  timeoutSeconds: 1800
```

The drill transitions through `Pending`, `Running`, and a terminal phase. The
status carries an ordered `evidence` list capturing each step the controller
executed and an aggregated `durationSeconds` field.

### Detecting Replica Skew

`RebalanceCheck` evaluates the topology distribution of a workload and surfaces
drift through standard conditions. It does not move pods. Mitigation, when
authorised by a parent `ResiliencePolicy`, is applied separately.

```yaml
apiVersion: resilience.hilios.io/v1alpha1
kind: RebalanceCheck
metadata:
  name: payments-zones
spec:
  targetSelector:
    matchLabels:
      app: payments
  topologyKey: topology.kubernetes.io/zone
  maxSkew: 1
  dryRun: true
```

### Reporting Contention

`ContentionReport` records noisy-neighbor signals and recommends mitigations
without applying them when `recommendOnly` is true.

```yaml
apiVersion: resilience.hilios.io/v1alpha1
kind: ContentionReport
metadata:
  name: payments-contention
spec:
  targetSelector:
    matchLabels:
      app: payments
  windowMinutes: 15
  thresholds:
    cpuStealPercent: 10
    memoryPressureEvents: 5
    throttlingPercent: 15
  recommendOnly: true
```

## API Reference

### ResiliencePolicy

| Field | Type | Description |
|---|---|---|
| `spec.targetSelector` | LabelSelector | Selects governed workloads (StatefulSets) |
| `spec.verifications[]` | VerificationSpec | Authorised verification kinds |
| `spec.mitigations[]` | MitigationKind | Authorised corrective actions |
| `spec.slo.recoveryTimeSeconds` | int32 | Maximum time for a recovery action |
| `spec.slo.maxReplicaSkew` | int32 | Maximum acceptable replica skew |
| `spec.schedule` | string | Optional cron expression |
| `spec.suspend` | bool | Pause reconciliation |

Conditions reported on status: `Ready`, `Validated`, `Degraded`.

### RecoveryDrill

| Field | Type | Description |
|---|---|---|
| `spec.type` | DrillType | `RestoreVerification` or `FailoverDrill` |
| `spec.policyRef` | LocalObjectReference | Optional parent policy |
| `spec.source` | BackupSource | Backup to restore (RestoreVerification) |
| `spec.verificationNamespace` | string | Temporary namespace name |
| `spec.healthChecks[]` | HealthCheck | Probes evaluated post-restore |
| `spec.cleanup` | bool | Delete artefacts after completion |
| `spec.timeoutSeconds` | int32 | Maximum drill duration |

Phases: `Pending`, `Running`, `Succeeded`, `Failed`.
Conditions: `Scheduled`, `Running`, `Succeeded`, `Failed`.

### RebalanceCheck

| Field | Type | Description |
|---|---|---|
| `spec.targetSelector` | LabelSelector | Pods inspected by the check |
| `spec.topologyKey` | string | Node label key (default `kubernetes.io/hostname`) |
| `spec.maxSkew` | int32 | Maximum acceptable skew |
| `spec.schedule` | string | Optional cron expression |
| `spec.dryRun` | bool | Report only, never recommend |

Conditions: `Balanced`, `Drifted`, `ActionRequired`.

### ContentionReport

| Field | Type | Description |
|---|---|---|
| `spec.targetSelector` | LabelSelector | Pods evaluated for contention |
| `spec.windowMinutes` | int32 | Look-back window for signals |
| `spec.thresholds` | ContentionThresholds | CPU steal, throttling, memory pressure |
| `spec.recommendOnly` | bool | Disable active mitigation |

Conditions: `Observed`, `Mitigated`, `Escalated`.

## Architecture

HILIOS is built with the [Operator SDK](https://sdk.operatorframework.io/) and
controller-runtime. It consists of four reconcilers backed by a shared internal
package set:

- `internal/conditions` wraps `metav1.Condition` so reconcilers update status
  through a small, typed surface.
- `internal/predicates` filters update events to generation changes and to
  toggling of the `hilios.io/paused` annotation, eliminating reconcile churn.
- `internal/finalizers` centralises the canonical finalizer string and helpers.
- `internal/topology` computes replica distribution and skew across topology
  domains and is fully unit-tested without a Kubernetes API.
- `internal/cronexpr` is a thin wrapper around `robfig/cron/v3` so scheduling
  has a single dependency surface.

Each reconciler:

1. Reads the resource and validates the spec.
2. Computes the observed state (matched workloads, distribution, findings).
3. Updates `.status.conditions` and a small set of typed status fields.
4. Emits Kubernetes Events for human-visible state transitions.
5. Requeues at a bounded interval for steady-state reconciliation.

## Features

- Declarative ResiliencePolicy with target selector, verifications, and mitigations
- RecoveryDrill state machine with audit-friendly evidence trail
- RebalanceCheck across arbitrary topology keys
- ContentionReport with pluggable signal evaluator (default: PodConditions)
- Standard Kubernetes API: conditions, events, finalizers, generation tracking
- Helm chart with RBAC, leader election, ServiceMonitor
- KUTTL end-to-end tests
- Dependabot with auto-merge for minor/patch updates
- govulncheck and gosec scans on every PR

## Development

See [DEVELOPER.md](DEVELOPER.md) for local development setup, including how to
run the operator against a kind cluster and how to run unit, integration, and
KUTTL tests.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

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
