# hilios-operator Helm Chart

Installs the HILIOS Operator into a dedicated namespace.

## Quickstart

```bash
helm install hilios . \
  --namespace hilios-system --create-namespace
```

## Values

See [values.yaml](values.yaml) for the full list. Key fields:

| Key | Default | Description |
|---|---|---|
| `replicaCount` | 1 | Number of manager replicas |
| `image.repository` | ghcr.io/kitio-tek/hilios-operator | Container image |
| `image.tag` | (chart appVersion) | Image tag |
| `leaderElection.enabled` | true | Enable leader election |
| `metricsPort` | 8443 | Prometheus metrics port |
| `healthProbePort` | 8081 | Liveness/readiness port |
| `crds.install` | true | Install CRDs from the chart |
| `rbac.create` | true | Create the cluster role binding |
| `serviceMonitor.enabled` | false | Create a Prometheus ServiceMonitor |
| `networkPolicy.enabled` | false | Install a NetworkPolicy for the manager |

## Compatibility

The chart targets Kubernetes 1.27 and later. The chart version tracks the
operator version.
