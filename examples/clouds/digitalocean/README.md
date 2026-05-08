# digitalocean examples

These manifests illustrate HILIOS resources tailored to digitalocean hosted
Kubernetes flavours. Replace placeholder selectors and topology keys with
values that match your actual workloads.

## Files

- `policy.yaml` - opt-in ResiliencePolicy for workloads on digitalocean.
- `rebalance.yaml` - zone-spread RebalanceCheck.
