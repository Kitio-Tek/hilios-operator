# ADR-0003 HILIOS does not inject failures

## Status

Accepted, 2026-05-08.

## Context

Litmus and similar projects intentionally inject failures to validate that a
workload survives. HILIOS sits in the same operational space but answers a
slightly different question: did the recovery path actually work the last
time we needed it.

## Decision

HILIOS controllers never injure a workload. Probes are read-only and
mitigations are recommendations except in narrow, well-bounded cases like
applying a topology spread constraint that is itself reversible.

## Consequences

Users do not have to grant HILIOS destructive RBAC. The trade-off is that
HILIOS cannot validate behaviours that only manifest under chaos; pair it
with a chaos tool when those scenarios matter.
