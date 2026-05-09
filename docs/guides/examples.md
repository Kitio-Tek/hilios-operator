# Examples

The `examples/` tree contains HILIOS manifests organised by intent. Use them
as starting points - each manifest is self-contained so `kubectl apply -f
<file>` succeeds against an empty cluster.

| Path | Contents |
|---|---|
| `examples/policy-basic.yaml` | Minimal ResiliencePolicy |
| `examples/policy-strict.yaml` | Stricter ResiliencePolicy with mitigations |
| `examples/drill-restore-verification.yaml` | RecoveryDrill referencing a Velero backup |
| `examples/rebalance-check.yaml` | Zone-spread RebalanceCheck |
| `examples/contention-report.yaml` | Recommend-only ContentionReport |
| `examples/catalog/` | Reusable templates per workload kind, per topology key, per posture |
| `examples/integration/` | Minimal StatefulSet shells with the opt-in label |

Each subdirectory has its own README with the file list.

## See also

- [README](../../README.md)
- [docs index](../README.md)
