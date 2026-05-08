# ADR-0004 Helm as the primary distribution channel

## Status

Accepted, 2026-05-08.

## Context

Operators are commonly distributed through OLM bundles, raw kustomize, or
Helm charts. Each has trade-offs. We need to pick a default for v0.x.

## Decision

Ship a Helm chart as the primary channel. OLM bundle support can be added
later without breaking existing installs.

## Consequences

Pros:

- Helm is widely available across managed Kubernetes platforms.
- The chart can express the deployment, RBAC, ServiceMonitor, NetworkPolicy
  in a single artefact.

Cons:

- OLM users have to wait for a future release.
- Helm hooks are required for ordering CRD installation; we lean on the
  built-in `crds/` directory in the chart instead.
