# Runbook: ResiliencePolicy reports Degraded

## Condition

`type: Degraded, status: True` on a ResiliencePolicy.

## Likely causes

- Selector matched no workloads.
- Most recent verification drill is older than the freshness window.
- Validation failed for a related resource.

## Diagnosis steps

1. `kubectl describe rpol <name>` to read events.
2. `kubectl get rdrill -A` to inspect drill history.
3. Check `.status.lastViolation` for a one-line summary.

## Remediation

- Run `kubectl apply` of a matching RecoveryDrill to reset the freshness window.
- If the selector is wrong, update `.spec.targetSelector` and re-apply.
