# Developer Guide

This document describes how to build, run, and test HILIOS Operator locally.

## Prerequisites

- Go 1.26 or later (matches the toolchain pinned in `go.mod` and CI)
- Docker
- kind 0.20 or later
- kubectl 1.27 or later
- Helm 3.x
- operator-sdk 1.38 or later
- KUTTL 0.15 or later (for end-to-end tests)
- Chainsaw 0.2.15 or later (for end-to-end tests)

## Repository Layout

```
api/v1alpha1/         CRD type definitions
cmd/                  Manager entrypoint
config/               Kustomize manifests produced by operator-sdk and kubebuilder
charts/hilios-operator/  Helm chart
internal/conditions/  metav1.Condition helpers
internal/controller/  Reconcilers, suite test, fake-client unit tests
internal/cronexpr/    Cron expression wrapper
internal/events/      Event recorder helpers
internal/finalizers/  Finalizer constants and helpers
internal/labels/      Well-known label constants
internal/predicates/  controller-runtime predicates
internal/topology/    Topology distribution and skew computation
tests/e2e/kuttl/      KUTTL test cases
tests/e2e/chainsaw/   Chainsaw test cases
test/e2e/             operator-sdk e2e scaffold
hack/                 Code generation boilerplate
```

## Building

```bash
make manifests generate
make build
```

`make manifests generate` regenerates the CRD YAML files in `config/crd/bases`
and the deepcopy methods in `api/v1alpha1/zz_generated.deepcopy.go`. The CI
pipeline fails the build if these files are out of date relative to the source
types.

## Unit Tests

```bash
go test ./...
```

Unit tests use the controller-runtime fake client and run in milliseconds. The
package `internal/controller` carries one test per controller covering the
happy path and the principal failure modes.

## Integration Tests with envtest

```bash
make envtest
KUBEBUILDER_ASSETS="$(setup-envtest use 1.30.0 -p path)" go test ./internal/controller/...
```

The Ginkgo suite in `internal/controller/suite_test.go` is skipped when envtest
binaries are not present, so the unit tests in the same package run cleanly in
isolation.

## Running Against a Local Cluster

Spin up a kind cluster and deploy the operator from the local working tree:

```bash
kind create cluster --name hilios
kubectl cluster-info --context kind-hilios

make docker-build IMG=hilios-operator:dev
kind load docker-image hilios-operator:dev --name hilios

helm install hilios charts/hilios-operator/ \
  --namespace hilios-system --create-namespace \
  --set image.repository=hilios-operator \
  --set image.tag=dev \
  --set image.pullPolicy=IfNotPresent
```

Apply a sample ResiliencePolicy:

```bash
kubectl apply -f - <<'EOF'
apiVersion: resilience.hilios.io/v1alpha1
kind: ResiliencePolicy
metadata:
  name: example
spec:
  targetSelector:
    matchLabels:
      hilios.io/enabled: "true"
  verifications:
    - kind: RestoreVerification
EOF
kubectl describe rpol example
```

Tear down:

```bash
kind delete cluster --name hilios
```

## End-to-End Tests

The repository ships two end-to-end suites that both run against a real
Kubernetes API. KUTTL is retained for the install path and the most basic
steady-state assertions; Chainsaw owns the new behavioural coverage where its
richer assertion language pays off.

### KUTTL

The KUTTL suite lives in `tests/e2e/kuttl`. Each numbered directory is a test
case; KUTTL applies the manifests in lexicographic order and verifies the
assertions.

```bash
helm install hilios charts/hilios-operator/ \
  --namespace hilios-system --create-namespace
make test-e2e-kuttl
# or directly
kubectl kuttl test tests/e2e/kuttl/ --config tests/e2e/kuttl/kuttl-test.yaml
```

### Chainsaw

The Chainsaw suite lives in `tests/e2e/chainsaw`. Each subdirectory contains a
`chainsaw-test.yaml` with `try`, `assert`, and `cleanup` blocks. Chainsaw is
preferred for tests that need:

- ordered or unordered list matching (`matchLabels` versus `initContainers`),
- comparative assertions (`length`, `>`, `<`, JMESPath),
- CLI or script verification alongside object assertions,
- per-step cleanup so the suite is hermetic when re-run against the same
  cluster.

```bash
helm install hilios charts/hilios-operator/ \
  --namespace hilios-system --create-namespace
make test-e2e-chainsaw
# or directly
chainsaw test tests/e2e/chainsaw/tests/ --config tests/e2e/chainsaw/.chainsaw.yaml
```

Install the Chainsaw CLI with `go install github.com/kyverno/chainsaw@latest`
or download the version pinned in `.github/workflows/e2e.yml`.

`make test-e2e` runs both suites sequentially.

## Linting

The repository uses `golangci-lint` with the configuration in `.golangci.yml`.
The CI pipeline pins golangci-lint to a known-good version (see
`.github/workflows/ci.yml`). To run the same checks locally:

```bash
golangci-lint run ./... --timeout 5m
```

## Releasing

Releases are tag-driven. Pushing a tag of the form `v<semver>` triggers the
`Release` workflow which builds and pushes the multi-arch container image to
GHCR, packages the Helm chart with the matching version, and creates a GitHub
Release with auto-generated notes.

```bash
git tag -a v0.1.0 -m "v0.1.0"
git push origin v0.1.0
```

## Troubleshooting

- `make manifests` fails with a `controller-gen` error: ensure `bin/controller-gen`
  is the version pinned in the `Makefile`. Removing `bin/controller-gen` forces
  a re-install on the next make invocation.
- envtest binaries missing: run `make envtest` or `setup-envtest use 1.30.0`
  and export `KUBEBUILDER_ASSETS` to the printed path.
- `helm template` succeeds but `helm install` fails on CRDs already present:
  pass `--skip-crds` (chart re-installation is non-destructive when `crds.keep=true`).
