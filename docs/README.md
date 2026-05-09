# HILIOS Documentation Index

## Reference

| Page | Topic |
|---|---|
| [architecture.md](architecture.md) | Reconciler topology and internal package map |
| [conditions.md](conditions.md) | Every condition the operator emits |
| [metrics.md](metrics.md) | Prometheus collectors |
| [probes.md](probes.md) | HTTP, Pod, Cmd, Kubernetes probe semantics |
| [comparison.md](comparison.md) | Relationship to Litmus, Velero, vendor operators |
| [roadmap.md](roadmap.md) | Planned releases |
| [quickstart.md](quickstart.md) | Kind cluster install walkthrough |
| [tools.md](tools.md) | Pinned development tool versions |

## Guides

| Page | Topic |
|---|---|
| [guides/examples.md](guides/examples.md) | Manifest examples by intent |
| [guides/observability.md](guides/observability.md) | Metrics, events, logs |
| [guides/scaling.md](guides/scaling.md) | Single-leader model and replica recommendations |
| [guides/quality.md](guides/quality.md) | Static analysis, security scanners, test runners |
| [guides/env-vars.md](guides/env-vars.md) | Environment variables honoured by the manager |
| [guides/release-checklist.md](guides/release-checklist.md) | Steps before tagging a release |
| [guides/cherrypick.md](guides/cherrypick.md) | Cherry-picking to release branches |

## Operations

| Page | Topic |
|---|---|
| [runbooks/policy-degraded.md](runbooks/policy-degraded.md) | ResiliencePolicy reports Degraded |
| [runbooks/drill-failed.md](runbooks/drill-failed.md) | RecoveryDrill ended in Failed phase |
| [runbooks/replica-skew.md](runbooks/replica-skew.md) | RebalanceCheck reports drift |
| [runbooks/contention-observed.md](runbooks/contention-observed.md) | ContentionReport observed signals |
| [runbooks/operator-not-ready.md](runbooks/operator-not-ready.md) | Manager pod not ready |

## Architectural Decision Records

See [adr/README.md](adr/README.md).
