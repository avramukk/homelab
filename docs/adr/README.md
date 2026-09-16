# Architecture Decisions (ADRs)

Every significant decision is recorded here as an **Architecture Decision
Record**. The goal is that a reader can reconstruct *why* the system looks the
way it does, including the alternatives that were rejected.

## Format

ADRs use a lightweight MADR-style layout:

```markdown
# ADR-NNN: <decision in one line>

## Status
Accepted | Planned | Superseded by ADR-NNN

## Date
YYYY-MM-DD

## Context
The forces at play — constraints, requirements, prior state.

## Decision
What we chose.

## Alternatives considered
Each option with pros / cons and why it was not chosen.

## Consequences
What becomes easier, what becomes harder, what we accept as a known limitation.
```

## Conventions

- **Numbering:** zero-padded, sequential, never reused.
- **Immutability:** an ADR is appended, never rewritten. A reversal is a new ADR
  with `Superseded by` / `Status: Superseded` cross-links.
- **One ADR per commit:** `docs(adr): ADR-NNN <short title>`.
- **Decisions that are expensive to reverse** always get an ADR.

## Index

| ADR | Title | Status |
|---|---|---|
| [001](001-runtime-talos-linux.md) | Cluster runtime: Talos Linux over Colima/k3s | Accepted |
| [002](002-cluster-topology.md) | Cluster topology: 1 control-plane + 2 workers | Accepted |
| [003](003-gitops-argocd.md) | GitOps controller: Argo CD | Accepted |
| [004](004-ci-github-actions.md) | CI: GitHub Actions | Accepted |
| [005](005-secrets-sealed-secrets.md) | Secrets in Git: Sealed Secrets | Accepted |
| [006](006-ingress-traefik.md) | Ingress: Traefik | Accepted |
| [007](007-observability-lgtm-otel.md) | Observability: self-hosted LGTM + OpenTelemetry | Accepted |
| [008](008-public-dashboard-cloudflare-tunnel.md) | Public dashboard via Cloudflare Tunnel | Accepted |
| [009](009-alerting-telegram-slo.md) | Alerting: Telegram, two severities, SLO burn-rate | Accepted |
| [010](010-backup-restic-local.md) | Backups: in-cluster Restic to local disk | Accepted |
| [011](011-iac-opentofu.md) | IaC: OpenTofu for external resources | Accepted |
| [012](012-security-baseline.md) | Security hardening baseline | Accepted |
| [013](013-docs-and-history-conventions.md) | Documentation and git-history conventions | Accepted |
| [014](014-incident-drill.md) | Deliberate incident drill and postmortem | Accepted |
| [015](015-repository-layout.md) | Monorepo layout | Accepted |
