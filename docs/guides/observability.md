# Observability

HILIOS exposes telemetry through three channels:

1. Prometheus metrics under the `hilios_*` prefix (see `docs/metrics.md`).
2. Kubernetes Events for human-visible state transitions (see
   `docs/conditions.md` for the reasons HILIOS emits).
3. controller-runtime structured logs via zap.

## Recommended scrape config

```yaml
- job_name: hilios-operator
  kubernetes_sd_configs:
    - role: endpoints
      namespaces:
        names: [hilios-system]
```
