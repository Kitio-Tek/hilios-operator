# Examples

Sample manifests for the HILIOS Operator CRDs.

| File | Resource | Notes |
|---|---|---|
| `policy-basic.yaml` | ResiliencePolicy | Minimal policy with a single restore verification |
| `policy-strict.yaml` | ResiliencePolicy | Three verifications, three mitigations, scheduled |
| `drill-restore-verification.yaml` | RecoveryDrill | Restore verification with Velero backup |
| `rebalance-check.yaml` | RebalanceCheck | Zone-spread check on `topology.kubernetes.io/zone` |
| `contention-report.yaml` | ContentionReport | Recommend-only contention evaluation |

Apply any of them after the operator is installed:

```bash
kubectl apply -f examples/policy-basic.yaml
kubectl describe rpol basic
```
