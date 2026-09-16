# ADR-004: CI — GitHub Actions

## Status
Accepted

## Date
2026-09-16

## Context

The repository needs continuous integration: validation (manifests, docs,
policies), documentation generation, and dependency/ image updates. The
repository is already hosted on GitHub and will be public.

Constraints:

- No self-hosted runner infrastructure to maintain.
- CI must not need cluster credentials (deployment is pull-based via
  [ADR-003](003-gitops-argocd.md)).
- Must run on a public repository without exposing secrets.

## Decision

Use **GitHub Actions** with hosted runners. CI scope, deliberately:

1. **Validate** — `markdownlint`, mermaid rendering, link checking.
2. **Validate manifests** — `kubeconform`, kustomize/Helm lint.
3. **Docs** — render/inspect generated docs where applicable.
4. **Updates** — Renovate PRs for images and dependencies (native Argo CD/Helm
   managers).

No push-deploy steps. Argo CD reconciles from `main`.

## Alternatives considered

### Self-hosted Git (GitLab CE)
- Pros: end-to-end self-hosting story; full DevSecOps platform.
- Cons: heavy for a three-node cluster (~8 GiB baseline); duplicates what GitHub
  already provides; resource cost crowds out the actual workloads.
- Rejected.

### Forgejo/Gitea + Woodpecker (or Gitea Actions)
- Pros: the homelab consensus; lightweight; dogfooding your own infrastructure.
- Cons: a second Git host alongside GitHub, which is already the public face;
  operational overhead for little showcase gain; runners would consume cluster
  resources.
- Rejected: the public repo lives on GitHub, so CI lives on GitHub too.

### Jenkins
- Pros: ubiquitous in enterprises.
- Cons: heavy, plugin drift, weak fit for a docs/validation workload.
- Rejected.

## Consequences

- **Gained:** zero runner maintenance, free minutes for a public repo, tight
  integration with Renovate, and a clean separation of concerns — CI verifies,
  Argo CD deploys.
- **Accepted:** CI is a cloud dependency; if GitHub is unavailable, validation
  pauses (deployment still reconciles independently).
- **Accepted:** secrets for CI, if any, are GitHub repository secrets — never
  cluster credentials.
- Auto-generated manifest PRs (CI opening config PRs) were considered and
  deferred; Renovate covers the image-update case for now.
