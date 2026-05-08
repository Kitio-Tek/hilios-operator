# Quickstart

This guide walks you through installing HILIOS Operator on a kind cluster,
applying a sample policy, and inspecting the operator's status output.

## Prerequisites

- Docker
- [kind](https://kind.sigs.k8s.io/) 0.22 or later
- Helm 3.x
- kubectl 1.27 or later

## 1. Create a Kind Cluster

```bash
kind create cluster --name hilios --image kindest/node:v1.30.0
kubectl cluster-info --context kind-hilios
```

## 2. Build and Load the Operator Image

```bash
make docker-build IMG=hilios-operator:dev
kind load docker-image hilios-operator:dev --name hilios
```

## 3. Install via Helm

```bash
make chart-sync
helm install hilios charts/hilios-operator/ \
  --namespace hilios-system --create-namespace \
  --set image.repository=hilios-operator \
  --set image.tag=dev \
  --set image.pullPolicy=IfNotPresent

kubectl wait --for=condition=Available --timeout=120s \
  deployment/hilios-hilios-operator -n hilios-system
```

## 4. Apply a Sample Policy

```bash
kubectl apply -f examples/policy-basic.yaml
kubectl get rpol -A
kubectl describe rpol basic
```

## 5. Run a Drill

```bash
kubectl apply -f examples/drill-restore-verification.yaml
watch -n 1 kubectl get rdrill nightly-restore
```

The drill cycles through `Pending`, `Running`, and finishes in `Succeeded` or
`Failed`. Inspect the evidence log:

```bash
kubectl get rdrill nightly-restore -o jsonpath='{.status.evidence}' | jq
```

## 6. Tear Down

```bash
helm uninstall hilios -n hilios-system
kind delete cluster --name hilios
```

## Troubleshooting

- The operator pod crash-loops with a permission error: install the chart with
  `rbac.create=true` (the default) or apply the cluster role manifests in
  `config/rbac` manually.
- A drill stays in `Pending`: the controller needs the finalizer to be applied
  before transitioning. Check `kubectl get rdrill <name> -o yaml | yq .metadata.finalizers`.
- Metrics endpoint refuses connections: the chart defaults to
  `--metrics-secure=false` for local development. For production, set
  `extraArgs: ["--metrics-secure=true"]` and provide TLS material.
