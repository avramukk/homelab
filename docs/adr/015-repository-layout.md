# ADR-015: Monorepo layout

## Status
Accepted

## Date
2026-09-16

## Context

The lab produces several kinds of artifact: documentation and ADRs, Kubernetes
manifests for GitOps, OpenTofu for external resources, and the source of the demo
service. These could live in one repository or several.

The overriding requirement ([ADR-013](013-docs-and-history-conventions.md)) is a
single, readable narrative: a reviewer should follow decisions → infrastructure →
workloads without hopping between repositories.

## Decision

Use a **monorepo** with a conventional layout:

```
.
├── docs/            # ADRs, architecture, runbooks, operations, security, incidents
├── infra/           # cluster bootstrap + GitOps root (Argo CD app-of-apps)
│   ├── bootstrap/   #   how Argo CD itself is installed
│   └── apps/        #   Argo CD Applications and their manifests
├── terraform/       # OpenTofu for external resources (Cloudflare, Tailscale, B2)
├── apps/            # source of the demo service
├── .github/         # CI workflows, PR template
├── CONTRIBUTING.md
└── CHANGELOG.md
```

Argo CD points at `infra/apps`; CI validates `docs/`, `infra/`, and `apps/`.

## Alternatives considered

### Two repositories (config + code)
- Pros: smaller clones; independent CI; the config repo can stay minimal.
- Cons: splits the narrative; a decision about the demo service and its deployment
  spans two histories; more plumbing for a solo operator.
- Rejected.

### Three+ repositories (docs / config / app / terraform)
- Pros: maximal separation of concerns.
- Cons: cross-repo coordination overhead; the "read the story in order" goal
  breaks immediately.
- Rejected.

### Let the structure emerge later
- Pros: no premature commitment.
- Cons: the history would show churn as things move; layout is cheap to decide
  now and expensive to reorganise after commits reference paths.
- Rejected.

## Consequences

- **Gained:** one linear history spanning decisions and implementation; Argo CD
  and CI both target a single source of truth; onboarding is one clone.
- **Accepted:** the repository grows; without path-based CI filters, unrelated CI
  would run wastefully — filters are used (only `docs/`, `infra/`, `apps/`,
  `terraform/` changes trigger the relevant jobs).
- **Accepted:** no per-directory ownership rules are needed for a solo operator;
  if collaborators join, `CODEOWNERS` is the natural next step.
- Path layout is now a stable contract referenced by ADRs and runbooks.
