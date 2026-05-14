# Operations

Operational reference for running the HILIOS manager in a cluster. For deeper
guidance see [guides/scaling.md](guides/scaling.md) (replica model),
[guides/observability.md](guides/observability.md) and [metrics.md](metrics.md)
(telemetry), and [guides/env-vars.md](guides/env-vars.md) (environment
variables).

## Logging

The manager uses zap-stdlib logging via controller-runtime. Set
`--zap-log-level=debug` to see verbose output.

## Leader Election

Leader election is enabled by default and uses the Kubernetes lease primitive
in `coordination.k8s.io/v1`.

## Image

The container image is published to `ghcr.io/kitio-tek/hilios-operator`.
Multi-arch (amd64, arm64) builds are produced on tag push.

## Ports

The manager binds 8443 for metrics and 8081 for health probes. The
`--metrics-secure` flag defaults to true in the binary; the bundled Helm
chart passes `--metrics-secure=false` so the metrics endpoint is plain HTTP
unless the deployment template is overridden.

## RBAC

The cluster role grants list/watch on the `resilience.hilios.io` group plus
read-only access to `apps/v1` StatefulSets and `core/v1` Pods/Nodes/Namespaces.

## High Availability

Run >= 2 replicas with leader election; the chart bundles `values-ha.yaml` as a
starting point.

## Telemetry

Prometheus metrics are exposed under `hilios_*`; ServiceMonitor objects can be
installed via the chart.

## Compatibility Matrix

Tested on Kubernetes 1.28-1.30. Newer versions may work but are not part of CI.
See [SUPPORTED_VERSIONS.md](../SUPPORTED_VERSIONS.md) for the support policy.

## kubectl Helper

A tiny `kubectl` plugin lives at `hack/kubectl-hilios`. Drop it into your `PATH`
and invoke `kubectl hilios get -A` to list every Hilios resource at once.

## See also

- [README](../README.md)
- [docs index](README.md)
