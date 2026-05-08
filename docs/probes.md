# Probes

HILIOS RecoveryDrill objects can carry a list of `healthChecks` (probes) that
the controller evaluates as part of the drill workflow. The probe schema is
inspired by the Litmus Chaos probe model so platform engineers familiar with
chaos pipelines can reuse mental models. Unlike Litmus, HILIOS probes are
read-only by design: they verify state but never inject failures.

## Probe Types

### `HTTP`

Issues an HTTP GET against `spec.url` and compares the response status code to
`spec.expectedStatusCode` (default 200). Useful for synthetic readiness
endpoints exposed by restored workloads.

```yaml
healthChecks:
  - name: app-readiness
    type: HTTP
    url: http://payments.payments-restore-verify.svc/healthz
    expectedStatusCode: 200
    timeoutSeconds: 30
```

### `Kubernetes`

Verifies that a Kubernetes object referenced by `spec.resource` exists.
Useful for restore drills where the success criterion is "the StatefulSet is
recreated in the verification namespace".

```yaml
healthChecks:
  - name: sts-exists
    type: Kubernetes
    resource:
      apiVersion: apps/v1
      kind: StatefulSet
      namespace: payments-restore-verify
      name: payments
```

### `Pod`

Records the pod selector for external evaluation. The controller does not
perform readiness checks today; pair with an HTTP or Kubernetes probe for
actionable verification.

```yaml
healthChecks:
  - name: pods-ready
    type: Pod
    podSelector:
      matchLabels:
        app: payments
```

### `Cmd`

Records a shell command and expected output for execution by an external
runner pod. The controller marks the probe as Skip in the evidence log; the
runner is responsible for executing the command and amending status.

```yaml
healthChecks:
  - name: smoke
    type: Cmd
    command: "psql -c 'SELECT count(*) FROM ledger'"
    expectedOutput: "1000"
```

## Evidence

Every probe attempt produces an `EvidenceRecord` on the drill's status:

| Field | Description |
|---|---|
| `step` | Probe identifier (`healthcheck:<name>`) |
| `time` | When the probe ran |
| `result` | `Pass`, `Fail`, or `Skip` |
| `message` | Human readable summary |

A drill with one failing probe transitions to `Failed` with reason
`HealthCheckFailed`. Probes execute sequentially; the first failure
short-circuits the drill.

## Implementation Notes

- The HTTP runner is the package-level variable `internal/healthcheck.httpRunner`
  protected by an `sync.RWMutex`. Tests override it through `SetHTTPRunner`.
- The Kubernetes runner uses an unstructured client to avoid taking a hard
  dependency on every API group HILIOS may need to inspect.
- Both runners honour `TimeoutSeconds`; values <= 0 default to 30 seconds.
