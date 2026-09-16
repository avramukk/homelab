# Terraform / OpenTofu

Infrastructure-as-code for resources **outside** Kubernetes — currently
**Cloudflare DNS only** ([ADR-011](../docs/adr/011-iac-opentofu.md),
[ADR-018](../docs/adr/018-cloudflare-tunnel-cli-opentofu-dns.md)). In-cluster
resources stay with Argo CD; nothing here duplicates them.

## What is managed

| Resource | Purpose |
|---|---|
| `cloudflare_dns_record` | proxied CNAMEs `<hostname>.avramukk.com` → `<tunnel_id>.cfargotunnel.com` |

The **Cloudflare Tunnel itself is created with the `cloudflared` CLI**, not by
OpenTofu — the account-scoped API token can read tunnels but not create them, and
granting account-level write scope was judged disproportionate ([ADR-018](../docs/adr/018-cloudflare-tunnel-cli-opentofu-dns.md)).

## Prerequisites

- OpenTofu `>= 1.9`, Cloudflare provider v5.
- A Cloudflare **API token** scoped to the `avramukk.com` zone with
  `Zone → DNS → Edit` and `Zone → Zone → Read`.
- A tunnel id from `cloudflared tunnel create homelab` (see the runbook below).

## Create the tunnel (once, on the host)

```bash
export PATH="/opt/homebrew/bin:$PATH"
cloudflared tunnel login                # only if no ~/.cloudflared/cert.pem exists
cloudflared tunnel create homelab       # prints the tunnel id
cloudflared tunnel list
```

## Run OpenTofu

```bash
export PATH="/opt/homebrew/bin:$PATH"
cd terraform
tofu init
tofu plan
tofu apply
```

## Secrets and state

- `terraform.tfvars` (token, `tunnel_id`) and all state files are **gitignored**.
- Copy `terraform.tfvars.example` → `terraform.tfvars` and fill it in, or export
  `TF_VAR_cloudflare_api_token` / `TF_VAR_tunnel_id` for a single run.
- `.terraform.lock.hcl` **is committed** for reproducible provider versions.

## Publishing a hostname

Add the subdomain to `public_hostnames` and apply:

```hcl
public_hostnames = ["homelab", "status", "grafana"]
```

Each entry creates a proxied CNAME. Traffic only reaches a service once
`cloudflared` runs in the cluster **and** has a matching ingress rule.
