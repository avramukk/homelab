---
name: cloudflare-tunnel
description: Expose this homelab's public dashboards through a Cloudflare Tunnel — create and route a tunnel, run cloudflared in-cluster, manage DNS, and debug 502/1033 errors. Use when publishing Grafana or the status page, working with cloudflared, tunnel tokens, ingress rules, or when the public dashboard is unreachable.
---

# Cloudflare Tunnel

Public access is **outbound-only**: `cloudflared` dials out to Cloudflare, so no
inbound ports are opened on the host ([ADR-008](../../docs/adr/008-public-dashboard-cloudflare-tunnel.md)).

## Model in this repo

- One tunnel, multiple hostnames routed to in-cluster Services.
- Tunnel token stored as a `SealedSecret`; `cloudflared` runs as a Deployment.
- DNS records are managed by OpenTofu ([ADR-011](../../docs/adr/011-iac-opentofu.md)).

## Create a tunnel (one-time, operator machine)

```bash
cloudflared tunnel login
cloudflared tunnel create homelab
cloudflared tunnel list
```

Route DNS (prefer OpenTofu for anything lasting):

```bash
cloudflared tunnel route dns homelab grafana.<domain>
cloudflared tunnel route dns homelab status.<domain>
```

Ingress config (in-cluster ConfigMap):

```yaml
tunnel: <TUNNEL_ID>
credentials-file: /etc/cloudflared/creds.json
ingress:
  - hostname: grafana.<domain>
    service: http://grafana.observability.svc.cluster.local:80
  - hostname: status.<domain>
    service: http://uptime-kuma.status.svc.cluster.local:80
  - service: http_status:404
```

## Run in-cluster

```bash
kubectl -n cloudflare get deploy,pods
kubectl -n cloudflare logs deploy/cloudflared --tail=100
# token-based (no credentials file):
kubectl -n cloudflare create secret generic tunnel-token --from-literal=TUNNEL_TOKEN=<token>  # do not commit
```

## Verify from outside

```bash
curl -sS -o /dev/null -w '%{http_code}\n' https://grafana.<domain>
dig +short grafana.<domain>
```

## Debug

| Symptom | Meaning | First checks |
|---|---|---|
| `502 Bad Gateway` | cloudflared reached, backend refused | Service name/port in ingress; is the pod up? |
| `1033` | tunnel not connected | `cloudflared` pod logs; token valid? |
| DNS wrong | record missing/stale | `dig`; OpenTofu state |

```bash
kubectl -n cloudflare logs deploy/cloudflared | grep -i error
kubectl -n observability get svc grafana
```

## Guardrails

- The tunnel token is a secret — SealedSecret only.
- Dashboards are anonymous read-only; never expose admin endpoints or datasources
  with write access.
