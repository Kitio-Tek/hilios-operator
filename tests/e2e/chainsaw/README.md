# Chainsaw End-to-End Tests

This directory contains the [Chainsaw](https://github.com/kyverno/chainsaw)
test suite for HILIOS Operator. Chainsaw is a declarative end-to-end testing
tool for Kubernetes that complements the [KUTTL](https://github.com/kudobuilder/kuttl)
suite under `tests/e2e/kuttl`.

## Why both Chainsaw and KUTTL

KUTTL covers the steady-state assertions and the initial install path. Chainsaw
extends coverage to scenarios that KUTTL struggles with:

- **Ordered and unordered list assertions** — Kubernetes resources use both,
  and Chainsaw lets each test pick the right semantics.
- **Comparative assertions** — Chainsaw supports `>`, `<`, `equalTo`, and JMESPath
  expressions for status fields, which KUTTL does not.
- **CLI / script verification** — Chainsaw can assert on the output of a
  command, not just on the resulting Kubernetes objects.
- **Step-level timeouts and retries** — granular control beyond the per-suite
  timeout that KUTTL exposes.
- **Cleanup as a first-class step** — each test declares its own teardown,
  which keeps the suite hermetic when run repeatedly against the same cluster.

The migration is incremental: KUTTL continues to guard the install path while
Chainsaw owns the new behavioural assertions. Both suites run on every PR.

## Layout

```
tests/e2e/chainsaw/
  .chainsaw.yaml                Global chainsaw configuration
  tests/
    operator-install/           Verifies the deployment is Available
    resilience-policy/          ResiliencePolicy validation and conditions
    recovery-drill/             RecoveryDrill scheduling and phase transitions
    rebalance-check/            RebalanceCheck balanced/drifted assertions
    contention-report/          ContentionReport status conditions
```

## Running locally

```bash
# Pre-requisite: kind cluster reachable via KUBECONFIG with the operator installed.
make test-e2e-chainsaw
```

Or directly:

```bash
chainsaw test tests/e2e/chainsaw/tests/ --config tests/e2e/chainsaw/.chainsaw.yaml
```

Install Chainsaw (release pinned in `.github/workflows/e2e.yml`):

```bash
go install github.com/kyverno/chainsaw@latest
```
