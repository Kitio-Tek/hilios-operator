# Tools

This document inventories the tools used during development. CI installs the
exact pinned versions; matching them locally avoids "works on my machine"
surprises.

| Tool | Version | Purpose |
|---|---|---|
| Go | 1.25.10 | language toolchain |
| controller-gen | v0.18.0 | CRD and DeepCopy generation |
| golangci-lint | v2.5.0 | static analysis |
| gosec | latest | security analyzer |
| govulncheck | latest | dependency CVE scanner |
| KUTTL | v0.15.0 | end-to-end test runner |
| kustomize | v5.4.2 | manifest composition |
| envtest | release-0.18 | controller-runtime test environment |
| Helm | 3.x | chart packaging |
| kind | 0.20+ | local kubernetes |
