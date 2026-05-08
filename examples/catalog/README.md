# Drill, Check, and Policy Catalog

Reusable resource templates organised by category. The catalog is inspired by
Litmus ChaosHub but produces verification and policy resources rather than
chaos experiments.

## Restore-verification drills

| Template | Workload |
|---|---|
| `restore-postgres.yaml` | PostgreSQL StatefulSet |
| `restore-mysql.yaml` | MySQL StatefulSet |
| `restore-mongodb.yaml` | MongoDB ReplicaSet |
| `restore-redis.yaml` | Redis master |
| `restore-elasticsearch.yaml` | Elasticsearch cluster |
| `restore-kafka.yaml` | Kafka StatefulSet |

## Failover drills

| Template | Workload |
|---|---|
| `failover-statefulset.yaml` | Any StatefulSet |

## Rebalance checks

| Template | Topology key |
|---|---|
| `rebalance/zone-spread.yaml` | topology.kubernetes.io/zone |
| `rebalance/host-spread.yaml` | kubernetes.io/hostname |
| `rebalance/region-spread.yaml` | topology.kubernetes.io/region |

## Contention reports

| Template | Signal |
|---|---|
| `contention/cpu-throttle.yaml` | CPU throttling |
| `contention/memory-pressure.yaml` | Memory pressure |
| `contention/noisy-neighbor.yaml` | Combined CPU steal and pressure |

## Resilience policies

| Template | Posture |
|---|---|
| `policy/conservative.yaml` | Daily restore verification only |
| `policy/balanced.yaml` | Restore + replica placement, with mitigation |
| `policy/aggressive.yaml` | All verifications, all mitigations, frequent |

Each manifest contains a placeholder backup reference (`REPLACE_WITH_BACKUP_NAME`).
Replace it with the name of a Velero `Backup` object in the `velero` namespace
before applying.
