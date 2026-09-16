# Terraform / OpenTofu

Infrastructure-as-code for resources **outside** Kubernetes — Cloudflare DNS and
the Cloudflare Tunnel ([ADR-008](../docs/adr/008-public-dashboard-cloudflare-tunnel.md),
[ADR-011](../docs/adr/011-iac-opentofu.md)). In-cluster resources stay with Argo CD;
nothing here duplicates them.

## What is managed

| Resource | Purpose |
|---|---|
| `cloudflare_zero_trust_tunnel_cloudflared` | the tunnel `cloudflared` in the cluster joins |
| `cloudflare_zero_trust_tunnel_cloudflared_config` | tunnel ingress rules (created only when hostnames are published) |
| `cloudflare_dns_record` | proxied CNAMEs for published hostnames |

Currently **no hostnames are published** — the tunnel and zone are prepared so
that DNS/ingress can be added when a service exists.

## Prerequisites

- OpenTofu `>= 1.9`, Cloudflare provider v5.
- A Cloudflare **API token** scoped to the `avramukk.com` zone with
  `Zone → DNS → Edit` and `Account → Cloudflare Tunnel → Edit`.

## Run

```bash
export PATH="/opt/homebrew/bin:$PATH"
cd terraform
tofu init
tofu plan
tofu apply
```

## Secrets and state

- `terraform.tfvars` (which holds the token) and all state files are **gitignored**.
- Copy `terraform.tfvars.example` → `terraform.tfvars` and fill in the token.
  Alternatively export `TF_VAR_cloudflare_api_token` for a single run.
- `.terraform.lock.hcl` **is committed** for reproducible provider versions.

## Publishing a hostname

Add an entry to `public_hostnames` (see the example), then `tofu plan && tofu apply`:

```hcl
public_hostnames = {
  "status"  = "http://uptime-kuma.status.svc.cluster.local:80"
  "grafana" = "http://grafana.observability.svc.cluster.local:80"
}
```

Each entry creates a proxied CNAME `<hostname>.avramukk.com` pointing at the
tunnel and a matching tunnel ingress rule. `cloudflared` must be running in the
cluster for traffic to reach the service.
