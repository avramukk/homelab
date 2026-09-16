# ADR-020: Grafana stays private; only the status page is public

## Status
Accepted — amends [ADR-008](008-public-dashboard-cloudflare-tunnel.md)

## Date
2026-09-16

## Context

[ADR-008](008-public-dashboard-cloudflare-tunnel.md) published an **anonymous,
read-only Grafana** alongside the public status page, on the reasoning that a
live dashboard is the most convincing part of a showcase.

On review that was rejected. Grafana is an **interactive, high-information
surface**: dashboards carry label values, host names, service topology, and
query results; the explore/datasource APIs are reachable next to the UI. The
showcase benefit of a live anonymous Grafana does not justify exposing
infrastructure detail to the internet, and it is the kind of surface that is hard
to audit for accidental leakage.

## Decision

- **Grafana is private.** It is reachable only over **Tailscale** (for now via
  `kubectl port-forward`; a Tailscale-served hostname can replace that later).
  Anonymous access is disabled.
- **Only the status page** (Uptime Kuma) is public: it is non-interactive and
  discloses a single fact per check (up/down).
- **`homelab.avramukk.com` is reserved**, not served: it resolves through the
  tunnel and the catch-all answers `404`, holding the name for a future landing page.
- OpenTofu now distinguishes `public_hostnames` (DNS + tunnel ingress) from
  `reserved_hostnames` (DNS only, no ingress).

## Alternatives considered

### Anonymous read-only Grafana (ADR-008 as written)
- Pros: the most visual artifact; zero-config for a visitor.
- Cons: interactive surface, broad infrastructure disclosure, hard to audit.
- Rejected.

### Grafana behind Cloudflare Access (Zero Trust)
- Pros: public URL, authenticated entry, audit trail.
- Cons: requires Cloudflare Access configuration and account-level Zero Trust
  permissions; access policy becomes another moving part for a single-operator lab.
- Deferred: viable if a public-but-authenticated dashboard is ever wanted.

### Tailscale-only Grafana (chosen)
- Pros: no public surface, no new dependencies, uses the network already required
  for operations.
- Cons: a visitor cannot click through dashboards; the public showcase relies on
  the status page, the repository, and screenshots.
- Accepted: the repository and status page carry the public story; dashboards are
  demonstrated on demand.

## Consequences

- **Gained:** the public attack surface is limited to a status page and a DNS
  name; Grafana needs no hardening for hostile traffic.
- **Accepted:** no live clickable dashboard in the public showcase. Screenshots
  and the documentation carry that role.
- **Terraform model:** `public_hostnames` vs `reserved_hostnames` makes the
  intent explicit and prevents accidentally publishing a service.
- **Reviewable:** publishing anything new is a one-line change in
  `terraform/terraform.tfvars` plus a `tofu plan` review.
