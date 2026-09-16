# ADR-018: Cloudflare Tunnel via cloudflared CLI; OpenTofu manages DNS only

## Status
Accepted

## Date
2026-09-16

## Context

[ADR-008](008-public-dashboard-cloudflare-tunnel.md) publishes the public
dashboards through a Cloudflare Tunnel, and [ADR-011](011-iac-opentofu.md) puts
external resources under OpenTofu. The initial plan was to create the tunnel in
OpenTofu with `cloudflare_zero_trust_tunnel_cloudflared`.

In practice the account-scoped API token could **read** tunnels
(`GET /accounts/{id}/cfd_tunnel` → 200) but not **create** them
(`POST … → 403 Authentication error`). The token had DNS edit and tunnel read but
not `Account → Cloudflare Tunnel → Edit`, and granting broad account-level
write scope purely to create one tunnel is a permission expansion we do not want.

There is already a valid `~/.cloudflared/cert.pem` (origin certificate) on the
host, authorising tunnel operations for this account.

## Decision

- **Create and own the tunnel with the `cloudflared` CLI**, authenticated by the
  existing **origin certificate** (`cloudflared tunnel create homelab`), not by an
  account-scoped API token.
- **OpenTofu manages DNS records only**: a proxied `CNAME <hostname>.<zone>` →
  `<tunnel_id>.cfargotunnel.com`. The tunnel id is an input variable.
- **Tunnel ingress** (which service a hostname reaches) is configured on the
  cloudflared side, not in OpenTofu.

## Alternatives considered

### Grant the token `Account → Cloudflare Tunnel → Edit`
- Pros: one tool (OpenTofu) manages tunnel + DNS; fully declarative.
- Cons: broadens account-level write scope for little gain; token still exposed to
  CI/automation contexts; the CLI already has a safer, narrower auth path.
- Rejected: disproportionate privilege for one resource.

### Create the tunnel in the Cloudflare dashboard
- Pros: no CLI, no token.
- Cons: a click-ops step outside Git and outside OpenTofu; harder to reproduce.
- Rejected: the CLI path is scriptable and auditable.

### Keep everything in OpenTofu but create tunnels per-environment with a scoped token
- Deferred: revisit if the account ever needs many tunnels managed as code.

## Consequences

- **Gained:** no account-write API token; tunnel creation uses the origin
  certificate; DNS stays reviewed, versioned, and reproducible in OpenTofu.
- **Accepted:** the tunnel itself is not in Terraform state — it is created
  out-of-band. The id is passed in via `tunnel_id`, and `cloudflared tunnel
  create` is documented in the runbook.
- **Split of truth:** OpenTofu answers "which hostnames exist"; cloudflared
  answers "where they route". Both are documented; the duplication is deliberate.
- **Secrets:** the tunnel credentials (`<id>.json`) and the tunnel token live on
  the host / in a `SealedSecret`, never in Git.
