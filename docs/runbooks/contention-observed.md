# Runbook: ContentionReport observed contention

## Condition

`type: Observed, status: True` on a ContentionReport.

## Likely causes

- Noisy neighbour on a shared node (CPU steal).
- Memory pressure from siblings.
- CPU throttling due to limits set too low.

## Diagnosis steps

1. Inspect `.status.findings[]` for the impacted pods.
2. Read the `recommendation` on each finding.
3. Cross-check with `kubectl top pod -A` and `kubectl describe node`.

## Remediation

- Bump CPU/memory requests on the impacted workload.
- Apply pod anti-affinity to spread the workload off the noisy node.
- Tag the workload with a higher PriorityClass.
