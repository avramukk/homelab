# ADR-003: GitOps controller — Argo CD

## Status
Accepted

## Date
2026-09-16

## Context

Cluster state should be declared in Git and reconciled automatically, so that the
repository is the source of truth and the history of the system is readable. A
GitOps controller is the mechanism.

Constraints:

- Single cluster, single operator.
- Must integrate cleanly with the chosen CI ([ADR-004](004-ci-github-actions.md))
  and secret handling ([ADR-005](005-secrets-sealed-secrets.md)).
- The choice is itself part of the showcase, so market relevance matters.

## Decision

Use **Argo CD**.

## Alternatives considered

### Flux
- Pros: modular, multi-source (Git/OCI/S3), lower resource footprint; used by
  reference homelabs (e.g. the Mischa/Talos model).
- Cons: smaller market share; less commonly named in job requirements; no
  built-in UI.
- Rejected: Argo CD's ubiquity is the stronger signal for a portfolio, and its UI
  supports the "walk me through it" story.

### Plain CI-driven deployment (kubectl apply from CI)
- Pros: no extra component.
- Cons: push-based; cluster credentials live in CI; no continuous reconciliation —
  drift is invisible.
- Rejected: contradicts the GitOps principle and weakens the audit story.

### Jenkins-based CD
- Pros: familiar in enterprises.
- Cons: heaviest to run and maintain; not GitOps; poor fit for a single host.
- Rejected.

## Consequences

- **Gained:** continuous reconciliation (drift is detected and surfaced), an
  app-of-apps bootstrap pattern, native Renovate integration for image updates,
  and the most widely recognised GitOps tool.
- **Accepted cost:** Argo CD is a sizeable component for a three-node cluster
  (see resource budget in [ADR-007](007-observability-lgtm-otel.md)).
- Secret values are needed at deploy time; this is why
  [ADR-005](005-secrets-sealed-secrets.md) selects an in-cluster decryption
  step that does not require Argo CD to hold plaintext.
- Deployment is **pull-based**: CI validates, Argo CD reconciles. No cluster
  credentials are stored in GitHub Actions.
