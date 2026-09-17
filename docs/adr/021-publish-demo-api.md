# ADR-021: Publish the demo API at demo.avramukk.com

## Status
Amended by [ADR-022](022-demo-api-tailnet-only.md) — amends [ADR-020](020-grafana-private.md)

## Date
2026-09-16

## Context

[ADR-020](020-grafana-private.md) reduced the public surface to the status page
only, and reserved `homelab.avramukk.com`. The demo service
(`apps/demo` + PostgreSQL) is now the SLO-bearing workload, and the operator
wants to reach it through their own domain rather than port-forwarding — with
everything still routed through Cloudflare.

The demo API is **not read-only**: it exposes `POST /api/items` (a write) and
`GET /api/error` (a guaranteed failure used for alert demos). Publishing it opens
a write path to its database and lets anyone drive its error rate and SLOs.

## Decision

Publish the demo API at **`demo.avramukk.com`** through the existing Cloudflare
Tunnel, managed in OpenTofu like every other hostname:

- DNS + tunnel ingress in `terraform/` (`public_hostnames`).
- Grafana stays **private** — ADR-020 is unchanged on that point.
- The **data is disposable**: the database holds demo rows only. No personal data,
  no secrets, nothing that must be true.
- Mitigations for the write endpoint (Cloudflare Access, rate limiting, method
  restriction) are **explicitly deferred**, not forgotten — see Consequences.

## Alternatives considered

### Cloudflare Access in front of the demo API
- Pros: authenticated entry; the write endpoint stops being public.
- Cons: more setup, and the point of this endpoint is a frictionless demo.
- Deferred: the first mitigation to add if abuse appears.

### Rate limit / method restriction at Traefik
- Pros: keeps the API public but bounds write/abuse; exercises the ingress layer.
- Cons: routing through Traefik instead of straight to the Service adds a hop that
  must be wired and tested.
- Deferred: pairs naturally with installing Cilium/Gateway API work (ADR-012).

### Keep it private (port-forward over Tailscale)
- Pros: no public write path at all.
- Cons: contradicts the explicit request to reach it through the domain.
- Rejected by the user's decision.

## Consequences

- **Gained:** the demo is reachable at a stable public URL, which is what makes
  RED metrics, traces, SLO burn-rate and the runbooks demonstrable end to end.
- **Accepted risk:** anyone can create items and trigger errors. The blast radius
  is the demo database; recovery is "delete rows" or recreate the volume.
- **Noise risk:** public traffic can fire the SLO alerts. That is acceptable (and
  even useful), but it means `DemoErrorBudget*` alerts can be triggered by
  strangers; the runbook explains how to tell real failures from driven ones.
- **Public surface is now:** `status.avramukk.com` (Uptime Kuma),
  `demo.avramukk.com` (demo API). `homelab.avramukk.com` stays reserved.
  Grafana remains private.
- **NetworkPolicy:** the demo pod accepts ingress from the `cloudflare`
  namespace (cloudflared), `observability` (Prometheus) and `traefik` only.
