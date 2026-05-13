# Contributing

Thank you for considering a contribution to HILIOS Operator. The project follows
the conventions used by mature open-source operators built on the Operator SDK
and kubebuilder.

## Development Setup

See [DEVELOPER.md](DEVELOPER.md) for the toolchain, repository layout, and the
local development workflow.

## Branching and Commits

- Work happens on feature branches off `main`. Branch names follow the form
  `feature/<short-summary>` or `fix/<short-summary>`.
- Commits are small, thematic, and use the imperative mood in the subject line.
  Group related changes; do not bundle unrelated refactors into a feature commit.
- Each commit must build and pass `go test ./...` in isolation.

## Pull Requests

- Open a draft PR early to gather feedback. Mark it ready for review when the
  diff is final and the description explains the why, not just the what.
- Every PR runs the CI workflow. Linting, unit tests, build, and helm checks
  must be green before review.
- Generated files (`config/crd/bases`, `api/.../zz_generated.deepcopy.go`) must
  be committed alongside the source changes that produced them. The CI
  pipeline fails the build if generated artefacts are out of date.
- Add or update end-to-end tests for behaviours that affect the controller's
  steady state. The repository ships two complementary suites:
  - [KUTTL](https://github.com/kudobuilder/kuttl) under `tests/e2e/kuttl`
    guards the install path and basic steady-state assertions.
  - [Chainsaw](https://github.com/kyverno/chainsaw) under `tests/e2e/chainsaw`
    covers behaviours where KUTTL is too limited, in particular ordered
    versus unordered list assertions, comparative operators (`>`, `<`,
    `length`), JMESPath status queries, CLI / script verification, and
    per-step cleanup. New behavioural coverage should generally be added to
    the Chainsaw suite.

## Code Style

- The code base avoids ornamental comments. Prefer code that names things
  clearly over a comment that paraphrases the implementation.
- Keep packages small and focused. The internal packages map one-to-one with a
  single concern (conditions, finalizers, predicates, topology).
- Reconcilers update status through the helpers in `internal/conditions` and
  emit events through `internal/events`. Do not import `corev1.EventTypeNormal`
  directly from controllers.

## Tests

- Every controller carries fake-client unit tests covering the happy path and
  at least one validation or transition failure mode.
- Pure helpers (cron parsing, topology skew, finalizer manipulation) ship with
  table tests.
- KUTTL and Chainsaw tests cover end-to-end scenarios that exercise the CRD
  lifecycle on a real Kubernetes API. Both suites run on every PR via the
  `E2E` workflow.

## Releasing

Maintainers cut releases by pushing a `v<semver>` tag. The `Release` workflow
builds and pushes the multi-arch image, packages the Helm chart, and creates
the GitHub Release notes.

## Reporting Issues

Open a GitHub issue with a clear reproducer when possible. For potential
security issues, follow the process in [SECURITY.md](SECURITY.md) instead.
