# Runbook: RebalanceCheck reports replica skew

## Condition

`type: Drifted, status: True` and/or `type: ActionRequired, status: True`
on a RebalanceCheck.

## Likely causes

- Pods have been scheduled to a single zone after a rolling restart.
- Node taints prevent the scheduler from spreading replicas.
- The workload's PodSpec lacks topology spread constraints.

## Diagnosis steps

1. Inspect `.status.distribution` to see the per-domain replica counts.
2. `kubectl get nodes -L topology.kubernetes.io/zone` to verify the topology key resolves.
3. `kubectl describe pod <pod>` to look for scheduling events.

## Remediation

- Apply a topology spread constraint to the workload.
- If a node is cordoned, consider draining it gracefully.
- Re-balance manually with `kubectl delete pod` to trigger a re-schedule.
