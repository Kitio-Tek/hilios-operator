# production namespace examples

Sample HILIOS resources scoped to the `production` namespace.

| File | Resource |
|---|---|
| `policy.yaml` | ResiliencePolicy |
| `drill.yaml` | RecoveryDrill |
| `rebalance.yaml` | RebalanceCheck |
| `contention.yaml` | ContentionReport |

Apply the bundle with:

```bash
kubectl apply -f examples/namespaces/production/
```
