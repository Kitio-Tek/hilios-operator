# Drill Catalog

Reusable RecoveryDrill templates for common workload shapes. The catalog is
inspired by the Litmus ChaosHub but produces verification drills rather than
chaos experiments.

| Template | Workload | Highlights |
|---|---|---|
| `restore-postgres.yaml` | PostgreSQL StatefulSet | Velero restore, StatefulSet existence, readiness probe |
| `restore-mysql.yaml` | MySQL StatefulSet | Velero restore, primary pod readiness |
| `restore-mongodb.yaml` | MongoDB ReplicaSet | Velero restore, StatefulSet existence verification |
| `restore-redis.yaml` | Redis master | Velero restore, exporter metrics endpoint |
| `restore-elasticsearch.yaml` | Elasticsearch cluster | Velero restore, cluster health endpoint |
| `restore-kafka.yaml` | Kafka StatefulSet | Velero restore, broker pod readiness |
| `failover-statefulset.yaml` | Any StatefulSet | Failover drill, Pod and Cmd probes |

Each manifest contains a placeholder backup reference (`REPLACE_WITH_BACKUP_NAME`).
Replace it with the name of a Velero `Backup` object in the `velero` namespace
before applying.
