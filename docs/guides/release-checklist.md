# Release Checklist

Use this checklist before tagging a release.

- [ ] CI on `main` is green for all workflows.
- [ ] CHANGELOG.md has a section for the new version.
- [ ] `version` and `appVersion` in the Helm chart match the upcoming tag.
- [ ] All Dependabot PRs are merged or queued.
- [ ] govulncheck reports no findings on Go 1.25.10.
- [ ] gosec reports no findings.
- [ ] KUTTL e2e suite passes locally on a kind cluster.
- [ ] Documentation guides for any new feature are merged.
- [ ] Tag is signed and pushed: `git tag -a v0.5.0 -m "v0.5.0" && git push origin v0.5.0`.
- [ ] Release workflow has produced the GHCR image and the GitHub Release notes.
- [ ] CHANGELOG.md is updated to add `## [Unreleased]` placeholder again.
