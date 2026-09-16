# ADR-017: Argo CD bootstrap — Helm once, then self-managed

## Status
Accepted

## Date
2026-09-16

## Context

[ADR-003](003-gitops-argocd.md) chose Argo CD as the GitOps controller, but a
controller cannot install itself through GitOps before it exists. Something has
to create Argo CD *outside* the reconciliation loop exactly once, after which the
cluster should converge to Git and stay there.

Requirements:

- The bootstrap must be **reproducible** and **idempotent** — re-running it must
  not break a healthy install.
- After bootstrap, Argo CD should **manage itself** from Git like any other app,
  so a manual `helm upgrade` is not the normal path.
- The method should use the fewest moving parts on a three-node lab cluster.

## Decision

Bootstrap Argo CD with the **official `argo/argo-cd` Helm chart**, installed by a
short idempotent script, then hand control to a **root Application (app-of-apps)**
that points at `infra/apps` and manages the rest — including Argo CD itself.

```
infra/
├── bootstrap/
│   ├── install-argocd.sh     # helm upgrade --install + apply root-app.yaml
│   ├── argocd-values.yaml    # resource limits, server.insecure for proxied TLS
│   └── root-app.yaml         # app-of-apps → infra/apps
└── apps/                     # every file is a child Application (recurse: true)
```

- `install-argocd.sh` is the **only** imperative step; it is safe to re-run.
- `root-app.yaml` uses `directory.recurse: true` with automated `prune` and
  `selfHeal`, so the cluster state converges to Git and drift is corrected.
- Argo CD's own manifests are moved under Git management as soon as the chart is
  adopted (self-management), rather than being upgraded by hand thereafter.

## Alternatives considered

### `kubectl apply` the upstream `install.yaml`
- Pros: no Helm dependency; a single manifest.
- Cons: no parameterisation, no release tracking, and upgrades become manual
  manifest diffs; CRDs and RBAC are harder to tune.
- Rejected: Helm gives versioned, values-driven installs for the same effort.

### Argo CD Operator / ApplicationSet-only bootstrap
- Pros: a CRD-driven lifecycle; ApplicationSets scale to many clusters.
- Cons: the operator is another component to run and understand; a single
  cluster does not benefit from ApplicationSets' fan-out.
- Rejected: disproportionate to a single-cluster lab.

### Commit rendered manifests and `kubectl apply -k`
- Pros: fully declarative in Git with no Helm at install time.
- Cons: rendered manifests are large, noisy in diffs, and must be regenerated on
  every upgrade — the opposite of reviewable history.
- Rejected.

## Consequences

- **Gained:** a one-command, re-runnable bootstrap; versioned chart with tuned
  resources; the cluster converges to Git after the first step.
- **Accepted:** `helm upgrade` *can* still be run manually, which would drift from
  Git-managed state. The rule is: bootstrap once, then change Argo CD through Git.
- **Access:** no ingress yet. The API/UI is reached with
  `kubectl -n argocd port-forward svc/argocd-server 8080:80` over Tailscale;
  public exposure is deferred ([ADR-008](008-public-dashboard-cloudflare-tunnel.md)).
- **TLS:** `server.insecure: true` is set because TLS is terminated by the future
  ingress/tunnel; the API is only reachable inside the cluster or over Tailscale.
