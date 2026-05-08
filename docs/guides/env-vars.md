# Environment Variables

The manager honours the following environment variables. Most have flag
equivalents; the env var is the recommended channel when running under Helm.

| Variable | Flag | Description |
|---|---|---|
| `WATCH_NAMESPACE` | n/a | Comma-separated namespaces to scope watches; empty means all |
| `LEADER_ELECTION` | `--leader-elect` | Enable leader election |
| `METRICS_ADDR` | `--metrics-bind-address` | Override the metrics bind address |
| `HEALTH_ADDR` | `--health-probe-bind-address` | Override the probe bind address |
| `LOG_LEVEL` | `--zap-log-level` | One of debug, info, warn, error |

The chart wires these via the deployment's env block; the user can extend the
list with `extraEnv` once that field lands in v0.6.

## See also

- [README](../../README.md)
- [docs index](../README.md)
