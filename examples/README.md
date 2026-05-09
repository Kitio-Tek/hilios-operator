# Examples

Top-level index of HILIOS manifests grouped by intent.

| Path | Intent |
|---|---|
| `policy-basic.yaml` | Minimal ResiliencePolicy |
| `policy-strict.yaml` | Stricter ResiliencePolicy with mitigations |
| `drill-restore-verification.yaml` | RecoveryDrill that references a Velero backup |
| `rebalance-check.yaml` | Zone-spread RebalanceCheck |
| `contention-report.yaml` | Recommend-only ContentionReport |
| `catalog/` | Reusable templates per workload kind, per topology key, per posture |
| `integration/` | Minimal StatefulSet shells with the opt-in label |

Start with the curated `catalog/` templates if you are unsure which file to
copy.
