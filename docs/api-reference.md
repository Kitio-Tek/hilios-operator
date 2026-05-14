# API Reference

Field tables for the four `resilience.hilios.io/v1alpha1` CRDs. For the full
catalogue of conditions and reasons each resource reports, see
[conditions.md](conditions.md).

## ResiliencePolicy

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

## RecoveryDrill

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

## RebalanceCheck

| Field | Type | Description |
|---|---|---|
| `spec.targetSelector` | LabelSelector | Pods inspected by the check |
| `spec.topologyKey` | string | Node label key (default `kubernetes.io/hostname`) |
| `spec.maxSkew` | int32 | Maximum acceptable skew |
| `spec.schedule` | string | Optional cron expression |
| `spec.dryRun` | bool | Report only, never recommend |

Conditions: `Balanced`, `Drifted`, `ActionRequired`.

## ContentionReport

| Field | Type | Description |
|---|---|---|
| `spec.targetSelector` | LabelSelector | Pods evaluated for contention |
| `spec.windowMinutes` | int32 | Look-back window for signals |
| `spec.thresholds` | ContentionThresholds | CPU steal, throttling, memory pressure |
| `spec.recommendOnly` | bool | Disable active mitigation |

Conditions: `Observed`, `Mitigated`, `Escalated`.

## See also

- [README](../README.md)
- [docs index](README.md)
- [conditions.md](conditions.md) - condition types and reasons in detail
