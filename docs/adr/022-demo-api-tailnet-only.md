# ADR-022: Demo API is tailnet-only (amends ADR-021)

## Status
Accepted — amends [ADR-021](021-publish-demo-api.md)

## Date
2026-09-17

## Context

[ADR-021](021-publish-demo-api.md) published the demo API at
`demo.avramukk.com` through the Cloudflare Tunnel and accepted an open write
endpoint as a documented risk. The network model was then split: only the
landing page and the status page stay public, and everything else is reachable
only over the tailnet ([ADR-020](020-grafana-private.md) and the Terraform
`private_hostnames` set). `demo` now resolves to the tailnet IP via a DNS-only
`A` record and is routed by Traefik — it is **not** in the tunnel.

The security audit found the documentation still claimed a public demo API.

## Decision

`demo.avramukk.com` is **tailnet-only**:

- DNS: **DNS-only A record → the tailnet IP** (not proxied, not a tunnel CNAME).
- Ingress: Traefik, restricted to the tailnet CIDR, with a Let's Encrypt cert.
- **Not** present in the Cloudflare Tunnel ingress; nothing public points at it.
- The unauthenticated write endpoint is therefore reachable only from devices on
  the tailnet; no rate limiting or Cloudflare Access is required.

## Alternatives considered

### Keep it public (ADR-021 as written)
- Rejected: the stated goal became "public by exception only" — landing + status.
  An open write endpoint adds risk with no showcase value.

### Public behind Cloudflare Access
- Rejected for now: extra moving parts for a demo; revisit if public access is
  ever wanted.

## Consequences

- The public attack surface is exactly two hostnames: `homelab` (landing) and
  `status` (Uptime Kuma).
- `terraform/terraform.tfvars` lists `demo` under `private_hostnames`; the
  landing page and this ADR were updated to match.
- Demo write abuse is no longer a public concern; data remains disposable.
