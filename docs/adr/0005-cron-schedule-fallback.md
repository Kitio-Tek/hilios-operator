# ADR-0005 Default requeue when spec.schedule is empty

## Status

Accepted, 2026-05-08.

## Context

HILIOS resources may declare a cron schedule for periodic re-evaluation. They
may also leave the field empty, in which case the controller should still
make forward progress.

## Decision

When `spec.schedule` is empty, requeue using a per-controller default
(5 minutes for ResiliencePolicy, 2 minutes for RebalanceCheck and
ContentionReport). Invalid cron strings fall back to the same default.

## Consequences

Operators who copy a sample manifest without filling in a schedule still see
the operator make progress. Authors who want predictable runs can set an
explicit cron string.
