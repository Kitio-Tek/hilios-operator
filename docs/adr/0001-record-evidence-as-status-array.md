# ADR-0001 Record drill evidence as a status array

## Status

Accepted, 2026-05-08.

## Context

RecoveryDrill needs to expose what the controller did during the drill so that
operators can audit the run after the fact. Two natural homes for that data
are Kubernetes Events and the drill's own status subresource.

## Decision

Use a `[]EvidenceRecord` array on `.status`. Events are still emitted for
human visibility but the canonical machine-readable record lives on status.

## Consequences

Pros:

- The evidence is durable for the lifetime of the resource (events expire
  after one hour by default).
- Tools can read it through the same status subresource they already use for
  conditions.
- The `evidence` slice serialises naturally as YAML.

Cons:

- Status payloads can grow large; we mitigate this by recording one entry per
  drill step rather than one per probe attempt.
- Updating a slice requires a full status update rather than a JSON patch.
