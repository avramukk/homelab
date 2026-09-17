# ADR-023: Host ports stay LAN-reachable (accepted risk)

## Status
Accepted — qualifies [ADR-020](020-grafana-private.md)

## Date
2026-09-17

## Context

ADR-020 declared Grafana "Tailscale-only". In practice the admin UIs answer on the
host's **LAN** address too, and an attempt to fix that failed:

- k3d publishes container ports **80/443** and Colima forwards them to the host on
  **all** interfaces; the macOS firewall permits the Colima SSH forwarder
  (`/usr/bin/ssh`), so the ports are open on the LAN as well as the tailnet.
- A Traefik `ipAllowList` (allow `100.64.0.0/10`) was tried and **reverted**: the
  k3d load-balancer proxy **NATs the client address**, so Traefik sees the proxy's
  address, not the client's — the rule rejected *everyone*, tailnet included.

What is actually reachable from the LAN: Grafana and Argo CD **login pages**, and
the demo API. Grafana/Argo require credentials; the demo API has no auth, but
reaching a specific vhost requires spoofing the `Host` header (`Host:
demo.avramukk.com` against the LAN IP) — it is not browsable by accident.

## Decision

**Do not close the host ports.** Document the LAN reachability as an accepted
risk. The domain design stays tailnet-only (`A` → tailnet IP); only the host's own
ports are additionally reachable inside the home network.

## Alternatives considered

### Traefik `ipAllowList`
- Tried; broke all access because the k3d proxy discards the source address.
- Rejected (reverted in `fce0960`).

### Host firewall (`pf`) rule blocking 80/443 on the LAN interface
- Works, but adds a privileged host-level component (LaunchDaemon) whose
  misconfiguration can break Colima port-forwarding entirely.
- Deferred: worth doing only if untrusted devices share the LAN.

### Recreate the cluster without published host ports, expose via `tailscale serve`
- The cleanest "tailnet-only by construction" option, but it rebuilds the cluster
  and resets in-cluster state (PVCs for the demo DB, Uptime Kuma, Grafana prefs).
- Deferred: revisit on the next cluster rebuild, or when moving to dedicated hardware.

### Put an auth token on the demo write endpoint
- Cheap, and removes the only unauthenticated surface.
- Deferred: optional hardening; the demo is disposable and tailnet-by-DNS.

## Consequences

- **Docs must not overclaim.** ADR-020's "Tailscale-only" means "not published to
  the internet and resolvable only on the tailnet"; it does **not** mean the host
  ports are closed to the LAN.
- **No new privileged component** on the host.
- **Revisit trigger:** another device or person on the LAN you do not control.
