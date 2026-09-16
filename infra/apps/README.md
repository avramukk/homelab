# Apps (app-of-apps)

Each file in this directory is a child **Argo CD `Application`**. The root
Application (`../bootstrap/root-app.yaml`) reconciles this directory with
`directory.recurse: true`, so adding a file here is how a new app is deployed.

## Convention

- One file per app, named after the app: `<app>.yaml`.
- Each child Application points at the path that holds its manifests (commonly
  `infra/<app>` or a Helm chart), not at this directory.
- Use `sync-wave` annotations only when ordering genuinely matters.
- Keep the desired state in Git; the cluster converges to it.

## Status

Empty at bootstrap — the first child applications land with the platform
components (ingress, secrets, observability). GitOps self-management of Argo CD
itself is added once the chart is adopted ([ADR-017](../../docs/adr/017-argocd-bootstrap.md)).
