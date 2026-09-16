# ADR-006: Ingress — Traefik

## Status
Accepted

## Date
2026-09-16

## Context

The cluster needs an ingress/edge layer to route public and internal HTTP(S)
traffic: Grafana and the status page publicly, internal services for the operator.
k3s ships Traefik by default, but the runtime is Talos
([ADR-001](001-runtime-talos-linux.md)), so the choice is open.

At the time of this decision, **ingress-nginx was retired (March 2026)** — no
further releases, fixes, or CVEs — which removed the historical default and made
this an active decision rather than a legacy default.

## Decision

Use **Traefik** as the ingress controller, with the **Gateway API** as the
preferred API for new resources.

## Alternatives considered

### ingress-nginx
- Pros: historically the default; enormous body of examples.
- Cons: **retired in March 2026**; migrating later is more expensive than choosing
  a supported controller now.
- Rejected: an unsupported ingress controller is not production-grade.

### Envoy Gateway (Gateway API)
- Pros: Gateway API–native, backed by Envoy; strong upstream momentum.
- Cons: heavier operational model; less homelab mindshare; fewer drop-in examples.
- Rejected: viable, but a steeper path for a single-host lab.

### Contour
- Pros: simple, Gateway API support.
- Cons: smaller community; less documentation for the patterns needed here.
- Rejected.

### Cilium ingress / Gateway API
- Pros: reuses Cilium, which is already present for network policy
  ([ADR-012](012-security-baseline.md)); fewer components.
- Cons: couples ingress to the CNI; troubleshooting becomes CNI-shaped; less
  familiar to reviewers.
- Rejected: keep the CNI and the edge layer as separate concerns.

## Consequences

- **Gained:** a supported, actively developed ingress controller with Gateway API
  leadership, a built-in dashboard, and a documented migration path from
  ingress-nginx.
- **Accepted:** Traefik's CRDs (`IngressRoute`, `Middleware`) are their own
  vocabulary to learn; the lab prefers Gateway API resources where supported to
  stay portable.
- TLS for internal services is handled with locally trusted certificates; public
  TLS terminates at Cloudflare ([ADR-008](008-public-dashboard-cloudflare-tunnel.md)).
