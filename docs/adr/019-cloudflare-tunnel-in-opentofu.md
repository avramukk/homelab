# ADR-019: Cloudflare Tunnel managed in OpenTofu

## Status
Accepted

## Date
2026-09-16

## Context

[ADR-018](018-cloudflare-tunnel-cli-opentofu-dns.md) moved tunnel creation to the
`cloudflared` CLI because the account-scoped API token could read tunnels but not
create them (`POST …/cfd_tunnel` → `403`). That was a **permission gap**, not a
tooling limitation: once the token gained `Account → Cloudflare Tunnel → Edit`,
a create probe returned `200` and a delete probe returned `200`.

Two sources of truth (tunnel in the CLI, DNS in OpenTofu) are worse than one,
and the tunnel was brand new and unused, so switching cost was minimal.

## Decision

Manage the **tunnel, its remote configuration, and DNS entirely in OpenTofu**:

- `cloudflare_zero_trust_tunnel_cloudflared` — `config_src = "cloudflare"`
  (remote-managed).
- `cloudflare_zero_trust_tunnel_cloudflared_config` — ingress rules; currently a
  catch-all `http_status:404` until services exist.
- `cloudflare_dns_record` — proxied CNAMEs to `<tunnel_id>.cfargotunnel.com`.

`cloudflared` in the cluster joins using the **tunnel token** (a `SealedSecret`),
so ingress changes do not require a cloudflared restart.

The CLI-created tunnel was deleted and recreated by OpenTofu so the tunnel has a
single owner. The user's pre-existing tunnel (`ldpi`) was left untouched.

## Alternatives considered

### Keep ADR-018 (CLI creates tunnel, OpenTofu does DNS)
- Pros: no account-level tunnel permission on the API token.
- Cons: two tools for one concern; tunnel not in state; harder to reproduce.
- Rejected: the permission is now granted and verified; one owner is better.

### Import the CLI-created tunnel instead of recreating it
- Pros: no delete/recreate.
- Cons: the tunnel was created with local config; changing to remote-managed
  config on an in-use tunnel is riskier than recreating an empty one.
- Rejected: the tunnel had no connections and no services.

## Consequences

- **Gained:** the tunnel, its routing config, and DNS are one reviewable,
  reproducible unit; `tofu plan` shows exactly what changes.
- **Accepted:** the API token now carries `Account → Cloudflare Tunnel → Edit`.
  This is real privilege; it is mitigated by a short expiration and by keeping the
  token out of Git (only `terraform.tfvars`, gitignored).
- **Known provider limitation:** `cloudflare_zero_trust_tunnel_cloudflared_config`
  **cannot be destroyed by Terraform** — apply logs a warning. Removing it means
  editing the tunnel in the Cloudflare dashboard. Documented in the runbook.
- **Supersedes** [ADR-018](018-cloudflare-tunnel-cli-opentofu-dns.md).
