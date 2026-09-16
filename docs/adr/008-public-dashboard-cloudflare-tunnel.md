# ADR-008: Public dashboard — Cloudflare Tunnel

## Status
Accepted

## Date
2026-09-16

## Context

The showcase is meant to be *seen*: a visitor should be able to open a live
dashboard and a status page. At the same time, the lab runs on a laptop behind a
home network — exposing inbound ports would be both unsafe and operationally
fragile (dynamic IP, no DDoS protection, home router as the perimeter).

Two things must be public (Grafana read-only, status page); everything else stays
private.

## Decision

Expose Grafana (read-only) and the status page through **Cloudflare Tunnel**
(`cloudflared`), on the operator's existing domain. The tunnel is outbound-only:
**no inbound ports are opened on the host**. Public access is **anonymous and
read-only**; all other services remain reachable only over Tailscale.

## Alternatives considered

### Tailscale Funnel
- Pros: no extra vendor; uses the network already in place; outbound-only.
- Cons: less natural for a public, memorable URL on own domain; less commonly seen
  in production front-door designs; identity model is oriented to the tailnet.
- Rejected: good fit, second choice; Cloudflare gives a proper public URL and edge.

### Port-forward + dynamic DNS
- Pros: no third party.
- Cons: opens the home network to the internet with no WAF/DDoS layer; requires
  managing TLS and dynamic DNS; unacceptable risk.
- Rejected.

### Static dashboard export via CI
- Pros: nothing is exposed live; a generated artifact is published.
- Cons: not live — a dashboard that lags by a pipeline run undermines the "real
  system" claim.
- Rejected as the primary path: it may later complement the tunnel for archival
  snapshots.

## Consequences

- **Gained:** a public HTTPS URL with no inbound attack surface; Cloudflare
  terminates TLS and absorbs volumetric traffic; the zero-trust option
  (Cloudflare Access) remains available if the anonymity decision is revisited.
- **Dependency:** availability of the public dashboard depends on Cloudflare and
  on `cloudflared` running in-cluster.
- **Secret:** the tunnel token is sensitive and is stored as a
  `SealedSecret` ([ADR-005](005-secrets-sealed-secrets.md)).
- **Data exposure:** dashboards are public and anonymous. They must contain no
  secrets, credentials, or sensitive internal identifiers; this is a review item
  in [ADR-012](012-security-baseline.md).
- DNS records for the tunnel hostnames are managed by OpenTofu
  ([ADR-011](011-iac-opentofu.md)).
