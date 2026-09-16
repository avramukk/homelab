# ADR-011: IaC — OpenTofu for external resources

## Status
Accepted

## Date
2026-09-16

## Context

Much of the lab is managed by GitOps (in-cluster resources) and shell scripts
(bootstrap). But some resources live **outside** the cluster and outside
Kubernetes' reach:

- DNS records for the tunnel hostnames (Cloudflare)
- Tailscale configuration
- (later) the offsite backup bucket (Backblaze B2)

These need a declarative, reviewable, state-tracked tool — the same "no
click-ops" standard as the rest of the lab.

## Decision

Use **OpenTofu** to manage external resources only (Cloudflare DNS, Tailscale,
and the future B2 bucket). In-cluster state stays with Argo CD; host bootstrap
stays with scripts.

State is stored encrypted; no secrets or state files are committed.

## Alternatives considered

### Terraform
- Pros: the most widely recognised name; largest provider/module ecosystem.
- Cons: BSL licence (source-available, not OSI-open); vendor direction under IBM
  is a known uncertainty for some users.
- Rejected: OpenTofu is a drop-in fork with the same HCL and providers, so the
  skill transfers, while the licence and governance fit an open showcase better.

### Pulumi
- Pros: real programming languages; good DX.
- Cons: adds a language runtime and a different mental model; the HCL skill is
  what most SRE job descriptions name.
- Rejected.

### Manual `curl`/dashboard changes, or scripts
- Pros: no tooling.
- Cons: not reviewable as a diff, no state, drift is invisible; DNS changes
  documented nowhere.
- Rejected: contradicts the lab's principles.

## Consequences

- **Gained:** external resources are versioned, reviewed in PRs, and show a plan
  before applying; an additional, market-relevant IaC skill; state encryption.
- **Accepted:** a state backend must exist and be secured. Backend credentials are
  secrets and never committed.
- **Scope discipline:** OpenTofu does **not** manage the cluster or workloads —
  duplicating Argo CD would create two sources of truth.
- `tofu plan` is required before any `apply`; plan output is part of the PR
  evidence.
