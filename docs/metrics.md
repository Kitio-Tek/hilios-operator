# Metrics Reference

The HILIOS manager exposes the standard controller-runtime metrics plus the
HILIOS-specific collectors documented below. All HILIOS metrics carry the
`hilios_` prefix.

## Counters

| Metric | Labels | Description |
|---|---|---|
| `hilios_policy_evaluations_total` | `namespace`, `name` | Total evaluations of each ResiliencePolicy |
| `hilios_drills_total` | `type`, `result` | Drills that reached a terminal phase, partitioned by type and outcome |

## Gauges

| Metric | Labels | Description |
|---|---|---|
| `hilios_policy_drift_count` | `namespace`, `name` | Violations observed in the most recent policy evaluation |
| `hilios_rebalance_skew` | `namespace`, `name` | Most recent skew observed by a RebalanceCheck |
| `hilios_contention_findings` | `namespace`, `name` | Findings produced by the most recent ContentionReport evaluation |

## Histograms

| Metric | Labels | Description |
|---|---|---|
| `hilios_drill_duration_seconds` | `type`, `result` | Wall-clock duration of completed drills, exponential buckets from 1s to ~4096s |

## Scraping

The Helm chart exposes the metrics on port 8443 by default. To expose them as a
Prometheus Operator `ServiceMonitor`, set `serviceMonitor.enabled=true`:

```bash
helm upgrade --install hilios charts/hilios-operator/ \
  --namespace hilios-system \
  --set serviceMonitor.enabled=true
```

The default chart enables HTTP for the metrics endpoint
(`--metrics-secure=false`) so a development scrape works without certificates.
For production, set `--metrics-secure=true` via `extraArgs` and provide the
trust material to the manager.
