# ADR-024: Status page served by Gatus

## Status
Accepted — supersedes the status-page role of [ADR-008](008-public-dashboard-cloudflare-tunnel.md)

## Date
2026-09-17

## Context

[ADR-008](008-public-dashboard-cloudflare-tunnel.md) chose **Uptime Kuma** for the
public status page and it runs at `status.avramukk.com`. Reviewing it for a
"like the big companies" status page exposed a problem:

- Uptime Kuma exposes **no REST API for administration** — `/api/auth/login`
  returns `Cannot POST`, `/api/monitors` and `/api/status-pages` fall back to the
  SPA. The only admin interface is the **Socket.IO** protocol.
- That means monitors and the status page can only be created by driving
  Socket.IO from a third-party client, with limited support for the newest
  feature (grouping monitors on a status page), and with version coupling.

The audit had already flagged Uptime Kuma's monitors as **click-ops**: not in Git,
not reproducible.

## Decision

Serve the status page with **Gatus** (`ghcr.io/twin/gatus:v5.36.0`), configured
entirely as code:

- Endpoints live in a `ConfigMap` (`infra/status/gatus/manifests/configmap.yaml`):
  grouped into **Public**, **Platform**, **Observability**, **Workloads**.
- Checks run **from inside the cluster**, so tailnet-only services (Grafana, Argo
  CD, Prometheus, Loki, Tempo, Alloy) are monitored without being exposed.
- The page is published through the Cloudflare Tunnel at `status.avramukk.com`.
- Adding a service to the status page is one block in the ConfigMap — a PR, not a
  UI session.

Uptime Kuma stays deployed but is removed from the public tunnel; its role as the
status page is retired.

## Alternatives considered

### Uptime Kuma + a Socket.IO seeding Job
- Keeps the incumbent tool, but chains the config to a third-party protocol
  client; grouping support on the status page is incomplete.
- Rejected: fragile, and it does not actually achieve "status page as code".

### Uptime Kuma configured by hand
- Fastest, but leaves the click-ops gap the audit found.
- Rejected.

## Consequences

- **Gained:** the status page *is* code; onboarding a service is a one-line PR;
  nice grouped uptime view with history (SQLite persistence on a PVC).
- **Accepted:** a second status tool ran alongside for a while; Uptime Kuma is now
  idle (kept for its UI-based ad-hoc checks).
- **Alerting is separate:** Gatus is a status page, not the alerting pipeline.
  SLO/burn-rate alerts stay in Alertmanager (ADR-009) — delivery still pending.
- **Image pinning:** the version is pinned in the manifest (`v5.36.0`).
