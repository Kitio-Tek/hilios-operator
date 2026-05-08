# Examples

The `examples/` tree contains HILIOS manifests organised by intent. Use them
as starting points - they are exhaustive on labels and selectors so a
`kubectl apply -f <file>` succeeds against an empty cluster.

| Path | Contents |
|---|---|
| `examples/policy-basic.yaml` | Minimal ResiliencePolicy |
| `examples/policy-strict.yaml` | Stricter ResiliencePolicy with mitigations |
| `examples/drill-restore-verification.yaml` | RecoveryDrill referencing a Velero backup |
| `examples/rebalance-check.yaml` | Zone-spread RebalanceCheck |
| `examples/contention-report.yaml` | Recommend-only ContentionReport |
| `examples/catalog/` | Reusable templates per workload kind |
| `examples/clouds/` | Cloud-flavoured starting points |
| `examples/integration/` | Per-workload tier and per-app samples |
| `examples/namespaces/` | Per-namespace bundles |

Each subdirectory has its own README with the file list.
