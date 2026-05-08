# ADR-0002 Conditions over phase strings

## Status

Accepted, 2026-05-08.

## Context

Some operators report state through a single `.status.phase` enum. Others
expose a slice of `metav1.Condition` entries. HILIOS resources transition
between several states that are not strictly mutually exclusive.

## Decision

Always expose conditions; expose `phase` only on `RecoveryDrill` where there
is a single workflow that genuinely has a primary state.

## Consequences

Conditions let HILIOS report `Validated=True, Ready=True, Degraded=True`
simultaneously, which is useful when a policy is otherwise healthy but
matched no workloads. A single phase string would force the controller to
choose between the two facts.
