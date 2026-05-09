# Conditions Reference

HILIOS resources expose their state through `metav1.Condition` entries on the
status subresource. This document catalogues every condition type and the
reasons that may accompany them.

## ResiliencePolicy

| Type | Reasons | Description |
|---|---|---|
| `Ready` | `Ready`, `Reconciling`, `Suspended`, `ValidationFailed` | Policy has been processed and the controller is tracking it |
| `Validated` | `Ready`, `ValidationFailed` | Spec passes validation |
| `Degraded` | `SelectorEmpty`, `Reconciling`, `Ready` | Policy detected a violation or matched no workloads |

## RecoveryDrill

| Type | Reasons | Description |
|---|---|---|
| `Scheduled` | `Scheduled` | Drill has been queued for execution |
| `Running` | `Started`, `Completed` | Drill is currently executing or has just transitioned out of Running |
| `Succeeded` | `RestoreVerified`, `Completed` | Drill produced evidence and reached the success terminal phase |
| `Failed` | `RestoreFailed`, `TimeoutExceeded`, `HealthCheckFailed` | Drill did not complete successfully |

## RebalanceCheck

| Type | Reasons | Description |
|---|---|---|
| `Balanced` | `ReplicasBalanced`, `ReplicaSkewDetected` | Topology distribution within tolerance |
| `Drifted` | `ReplicaSkewDetected`, `ReplicasBalanced` | Topology distribution exceeds tolerance |
| `ActionRequired` | `ReplicaSkewDetected`, `ReplicasBalanced` | Recommendation for the operator |

## ContentionReport

| Type | Reasons | Description |
|---|---|---|
| `Observed` | `ContentionDetected`, `Ready` | Contention signals were observed in the look-back window |
| `Mitigated` | `MitigationDisallowed`, `Ready` | Always False today: HILIOS records recommendations on findings; active mitigation is on the roadmap |
| `Escalated` | `Escalated` | Reserved for future use: contention persisted past the policy threshold |

## Generic Reasons

These reasons appear across multiple condition types:

| Reason | Meaning |
|---|---|
| `Reconciling` | Controller is in the process of evaluating the resource |
| `ValidationFailed` | The spec failed validation |
| `Suspended` | Reconciliation is suspended via `spec.suspend` or annotation |
| `TimeoutExceeded` | A bounded operation exceeded its timeout |
| `Completed` | Generic completion marker for transitional state |
