# Code Quality

HILIOS uses several scanners and linters to keep the code clean and safe.

| Tool | Purpose | Workflow |
|---|---|---|
| golangci-lint v2.5.0 | Static analysis | `.github/workflows/ci.yml` |
| go vet | Stdlib analysis | `.github/workflows/ci.yml` |
| gofmt | Formatting | `.github/workflows/ci.yml` |
| gosec | Security analysis | `.github/workflows/security.yml` |
| govulncheck | Stdlib + module CVE scan | `.github/workflows/security.yml` |
| Trivy | Filesystem and image scan | `.github/workflows/trivy.yml` |
| CodeQL | Semantic vulnerability scan | `.github/workflows/codeql.yml` |
| Helm lint | Chart linting | `.github/workflows/helm.yml` |
| KUTTL | End-to-end tests | local + manual workflow |

The lint, security, and helm jobs are required for `main`.
