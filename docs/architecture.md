# Architecture

This document describes the internal structure of HILIOS Operator and the
boundaries between the four reconcilers.

## Reconciler Topology

```
                                  +-----------------------+
                                  |   ResiliencePolicy    |
                                  |       (parent)        |
                                  +-----------+-----------+
                                              |
            +---------------------------------+---------------------------------+
            |                                 |                                 |
            v                                 v                                 v
+-----------+-----------+         +-----------+-----------+         +-----------+-----------+
|     RecoveryDrill     |         |    RebalanceCheck     |         |   ContentionReport    |
|       (workflow)      |         |     (observation)     |         |     (observation)     |
+-----------------------+         +-----------------------+         +-----------------------+
```

`ResiliencePolicy` is the only parent CRD. The other three CRDs reference it
through `spec.policyRef`. The policy controller does not create child resources
implicitly. Drills, checks, and reports are created by users or higher-level
controllers (CronJobs, GitOps pipelines) and reference the parent policy when
attribution is required.

## Reconcile Loop

Every reconciler follows the same shape:

1. Fetch the resource. Return on `NotFound`.
2. Handle the deletion path if `metadata.deletionTimestamp` is set.
3. Add the canonical finalizer (`hilios.io/finalizer`) if missing and requeue.
4. Validate the spec. Update `Validated` and `Ready` conditions on failure.
5. Compute the observed state from the cluster (matched workloads, distribution,
   findings).
6. Update the status subresource through the controller-runtime client. On
   `Conflict` the reconcile returns an error and the manager requeues with
   exponential backoff. The `internal/statuswriter` helper is available for
   reconcilers that need an in-call retry loop; it is not yet wired into the
   built-in controllers.
7. Emit Kubernetes Events through `internal/events` for transitions a human
   operator should see.
8. Increment metrics in `internal/metrics`.
9. Compute the next requeue from `internal/scheduling` (cron-driven when
   `spec.schedule` is set, otherwise a per-controller default).

## Internal Packages

| Package | Responsibility |
|---|---|
| `internal/conditions` | Set, query, and remove `metav1.Condition` entries |
| `internal/cronexpr` | Parse five-field cron expressions |
| `internal/events` | Wrap controller-runtime event recorder |
| `internal/finalizers` | Canonical finalizer string + helpers |
| `internal/labels` | Well-known labels and annotations |
| `internal/metrics` | Prometheus collectors registered with the manager |
| `internal/predicates` | Filter reconcile events to generation/pause changes |
| `internal/scheduling` | Compute next requeue duration |
| `internal/statuswriter` | Retry-on-conflict status update helper (opt-in) |
| `internal/topology` | Compute replica distribution and skew |

Each package is small (typically under 100 lines) and ships with a focused unit
test that exercises the boundary it owns. None of these packages depend on
controller-runtime types except where explicitly required.

## Failure Modes

- **Validation failure**: `Validated` condition becomes False with reason
  `ValidationFailed`. `Ready` becomes False. No mitigations are applied.
- **API server conflict**: a status update conflict surfaces as a reconcile
  error and triggers controller-runtime exponential backoff. Reconcilers that
  opt into the `internal/statuswriter` helper retry in-call before backing off.
- **Selector matched no workloads**: `Degraded` condition becomes True with
  reason `SelectorEmpty`. The policy is still considered Ready so that other
  controllers depending on the policy can proceed.
- **Drill timeout**: `Failed` condition becomes True with reason
  `TimeoutExceeded`. Evidence captures the timeout step.

## Extensibility

The signal evaluator inside `ContentionReport` is intentionally pluggable:
replace `evaluatePod` to read from a Prometheus query, the metrics-server, or a
custom collector without changing any other controller code.

The drill executor inside `RecoveryDrill` follows the same contract:
`executeDrill` is the single function that performs the verification work and
can be swapped for a Velero-aware implementation, a script runner, or a custom
HTTP probe service.
