# Runbook: HILIOS manager pod not ready

## Symptom

`kubectl get pods -n hilios-system` shows the manager pod stuck in
`CrashLoopBackOff` or `ContainerCreating`.

## Likely causes

- Image pull failure (private registry credentials missing).
- RBAC denied: the manager cannot list its CRDs.
- Leader election lease unavailable.

## Diagnosis

- `kubectl logs -n hilios-system <pod>` for crash details.
- `kubectl describe pod <pod>` for image pull events.

## Remediation

- Provide an image pull secret via `helm install --set imagePullSecrets`.
- Ensure the chart was installed with `rbac.create=true`.
- For lease issues, redeploy with `leaderElection.enabled=false` if running a
  single replica.
