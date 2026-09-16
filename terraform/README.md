# Terraform / OpenTofu

Infrastructure-as-code for resources **outside** Kubernetes — Cloudflare DNS
**and** the Cloudflare Tunnel ([ADR-011](../docs/adr/011-iac-opentofu.md),
[ADR-019](../docs/adr/019-cloudflare-tunnel-in-opentofu.md)). In-cluster
resources stay with Argo CD; nothing here duplicates them.

## What is managed

| Resource | Purpose |
|---|---|
| `cloudflare_zero_trust_tunnel_cloudflared` | the tunnel (`homelab`), remote-managed config |
| `cloudflare_zero_trust_tunnel_cloudflared_config` | tunnel ingress rules (catch-all 404 until services exist) |
| `cloudflare_dns_record` | proxied CNAMEs `<hostname>.<zone>` → `<tunnel_id>.cfargotunnel.com` |

`cloudflared` in the cluster joins the tunnel with its **token** (stored as a
`SealedSecret`), so routing changes are applied from Cloudflare without restarting
the connector.

## Prerequisites

- OpenTofu `>= 1.9`, Cloudflare provider v5.
- A Cloudflare **API token** for `avramukk.com` with:
  - `Zone → DNS → Edit`
  - `Zone → Zone → Read`
  - `Account → Cloudflare Tunnel → Edit`

## Run

```bash
export PATH="/opt/homebrew/bin:$PATH"
cd terraform
tofu init
tofu plan
tofu apply
```

## Secrets and state

- `terraform.tfvars` (token) and all state files are **gitignored**.
- Copy `terraform.tfvars.example` → `terraform.tfvars`, or export
  `TF_VAR_cloudflare_api_token` for a single run.
- `.terraform.lock.hcl` **is committed** for reproducible provider versions.

## Publishing a hostname

Add an entry to `public_hostnames` (subdomain → in-cluster Service URL) and apply:

```hcl
public_hostnames = {
  "status" = "http://uptime-kuma.status.svc.cluster.local:3001"
}
```

Each entry creates a proxied CNAME **and** a matching tunnel ingress rule. The
service URL must include the **actual Service port** (e.g. Uptime Kuma listens on
`3001`, not `80`). `cloudflared` must be running in the cluster.

To **reserve** a name without serving anything (the tunnel catch-all answers
404), add it to `reserved_hostnames` instead:

```hcl
reserved_hostnames = ["homelab"]
```

Only put a service in `public_hostnames` when it is meant to be reachable by
anyone — Grafana is deliberately private ([ADR-020](../docs/adr/020-grafana-private.md)).

## Known limitations

- `cloudflare_zero_trust_tunnel_cloudflared_config` **cannot be destroyed by
  Terraform** (provider behaviour); removing ingress means editing the tunnel in
  the Cloudflare dashboard. `tofu apply` prints this warning.
