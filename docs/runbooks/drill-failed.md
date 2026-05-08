# Runbook: RecoveryDrill phase Failed

## Condition

`.status.phase = Failed` on a RecoveryDrill.

## Likely causes

- Backup source not reachable.
- Health check probe failed against the restored workload.
- Drill exceeded its `spec.timeoutSeconds`.

## Diagnosis steps

1. `kubectl get rdrill <name> -o yaml` and inspect `.status.evidence`.
2. Look for the `Fail` evidence record - it contains the immediate error.
3. `kubectl get events -n <ns>` to see the matching Warning event.

## Remediation

- For probe failures: tighten or relax `expectedStatusCode` on the probe.
- For timeouts: bump `spec.timeoutSeconds`.
- For source errors: confirm Velero is healthy via `kubectl get backup`.
