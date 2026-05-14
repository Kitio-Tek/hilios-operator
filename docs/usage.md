# Usage

HILIOS is driven entirely through its four CRDs. This page shows a minimal
manifest for each and the commands to inspect it. See
[api-reference.md](api-reference.md) for the full field list and
[guides/examples.md](guides/examples.md) for the manifests shipped under
`examples/`.

## Declaring a ResiliencePolicy

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

## Running a RecoveryDrill

A `RecoveryDrill` is a one-shot or scheduled verification. The current build
covers two drill types:

- **`RestoreVerification`** creates a verification namespace and runs the
  configured probes against it. The actual restore must be triggered
  externally - for example by a Velero scheduled restore that targets the
  verification namespace. Native Velero invocation is on the v0.6 roadmap.
- **`FailoverDrill`** is a placeholder type that runs the configured probes
  but does not yet drive the failover sequence itself. It is intended for
  external runners (Argo Workflows, CronJobs) that orchestrate the failover
  and rely on Hilios for evidence collection.

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

## Detecting Replica Skew

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

## Reporting Contention

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

## See also

- [README](../README.md)
- [docs index](README.md)
