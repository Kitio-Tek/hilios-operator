# Comparison with Adjacent Projects

HILIOS sits at the intersection of **resilience verification**,
**chaos engineering**, and **operator-driven control planes**. This document
explains how the project relates to a few common reference points.

## Litmus Chaos

[Litmus](https://github.com/litmuschaos/chaos-operator) is a Kubernetes
operator for **injecting** chaos into workloads. Its core CRDs (`ChaosEngine`,
`ChaosExperiment`) describe the experiment to run and the application under
test. HILIOS borrows the **probe vocabulary** (HTTP, Kubernetes, Cmd, Pod) but
inverts the use case:

| Aspect | Litmus | HILIOS |
|---|---|---|
| Posture | Inject failures | Verify readiness and policy |
| Primary CRD | `ChaosEngine` | `ResiliencePolicy` |
| Secondary CRD | `ChaosResult` | `RecoveryDrill`, `RebalanceCheck`, `ContentionReport` |
| Probe semantics | Read-only check during chaos | Read-only check during drill |
| Side effects | Yes (controlled failure injection) | No (read-only by design) |

If you need to *inject* failures, use Litmus. If you need to *prove the
recovery path works* and *enforce policy*, use HILIOS.

## CloudNativePG / Vertica Operators

Vendor operators like [CloudNativePG](https://github.com/cloudnative-pg/cloudnative-pg)
or [Vertica Kubernetes](https://github.com/vertica/vertica-kubernetes) own the
full lifecycle of a single workload kind. They install the workload, manage
its replicas, take its backups, and route its traffic.

HILIOS does **not** own the workload lifecycle. It runs alongside vendor
operators and validates the operational invariants those operators rely on:
backup freshness, restore correctness, replica spread, contention. The
mental separation:

- The vendor operator answers "is the cluster running?".
- HILIOS answers "would we survive if the cluster failed?".

## Velero

[Velero](https://github.com/vmware-tanzu/velero) provides Backup and Restore
CRDs. HILIOS treats Velero as a backend: a `RecoveryDrill` of type
`RestoreVerification` references a `Backup` object via `spec.source` and
expects Velero to perform the actual restore. HILIOS contributes the
verification harness (temporary namespace, probes, evidence log) on top of
Velero's primitives.

## Kubernetes Operators in General

HILIOS follows the
[Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) and
[Operator SDK](https://sdk.operatorframework.io/) conventions:

- Reconcilers are status-driven.
- All long-running state is exposed as `metav1.Condition` entries.
- Events accompany every important transition.
- Predicates suppress reconcile churn on noisy fields.
- Generated manifests are committed alongside the source types.

If you have built another operator with these tools, HILIOS will look
familiar within the first few minutes of reading `internal/controller`.
