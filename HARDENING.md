# Hardening Guide

The default Helm chart is opinionated for safety. Production deployments may
want to go further:

- Set `--metrics-secure=true` and provide TLS material via volume mounts.
- Apply `networkPolicy.enabled=true` and limit ingress to your Prometheus ns.
- Run `helm install ... -f charts/hilios-operator/values-strict.yaml` for the
  recommended defaults.
- Restrict the default ClusterRole to the workloads HILIOS actually governs
  using a Role + RoleBinding under each watched namespace.
- Pin the manager image by digest, not tag.
