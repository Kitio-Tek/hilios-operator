# Security Policy

## Supported Versions

The latest minor release receives security updates. Older minor versions are
patched only when the fix is non-invasive.

## Reporting a Vulnerability

Please report suspected vulnerabilities privately by opening a GitHub Security
Advisory on this repository. Do not file a public issue for security reports.

A maintainer will acknowledge the report within five working days. If the
report is confirmed, the fix is developed in a private fork and released
together with a CVE entry where applicable.

## Scope

In scope:

- The HILIOS Operator container image and its Helm chart.
- The `resilience.hilios.io` API surface and its custom resources.
- The default RBAC manifests shipped with the chart.

Out of scope:

- Third-party dependencies (please report upstream).
- Misconfiguration of the cluster on which HILIOS runs.

## Hardening Notes

- The default Helm chart runs the manager as a non-root user with a read-only
  root filesystem and `seccompProfile: RuntimeDefault`.
- HTTP/2 is disabled on the metrics and webhook servers by default to avoid
  the Stream Cancellation and Rapid Reset CVEs.
- The manager talks to the Kubernetes API only with the verbs declared in
  `charts/hilios-operator/templates/clusterrole.yaml`. Do not grant additional
  cluster-wide permissions.
