# Integration examples

Sample HILIOS manifests scoped to common workloads and tiers. The directory
holds:

- StatefulSet shells with the `hilios.io/enabled` label.
- ResiliencePolicy / RecoveryDrill / RebalanceCheck / ContentionReport per
  workload tier (`critical`, `high`, `medium`, `low`).
- Sample drill, rebalance, and contention CRs in the `sample-N`, `drill-N`
  series.

Apply them piecemeal as starting points; do not blanket-apply the directory
in production.
