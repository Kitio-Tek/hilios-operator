# Scaling

The HILIOS manager is single-leader. Scaling out replicas does not increase
throughput; it increases availability through faster leader election after a
crash.

## Recommendations

- Run 1 replica for development clusters.
- Run 2-3 replicas in production with leader election enabled.
- Pin the manager pods to nodes labelled `node-role.kubernetes.io/control-plane=`
  if your cluster has dedicated control-plane nodes.
