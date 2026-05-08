# Per-namespace examples

Each subdirectory contains the four HILIOS CRs scoped to a single namespace.
This is the layout we recommend for multi-tenant clusters where a separate
namespace owns each environment.

| Namespace | Use case |
|---|---|
| production | Production workloads |
| staging | Pre-production smoke tests |
| qa | Continuous integration test workloads |
| dev | Developer namespaces |
| sandbox | Throwaway experimentation |
| infra | Cluster infrastructure |
| observability | Prometheus, Grafana, Tempo, Loki |
| data | ETL pipelines |
| ml | ML training |
| ai | Inference serving |
