# Runbook: Cloudflare Tunnel / public dashboard

**Means:** a public hostname (e.g. `homelab.avramukk.com`) returns a Cloudflare
error (`1033`, `502`, `530`) or nothing resolves.
**First check:** `curl -s -o /dev/null -w '%{http_code}\n' https://homelab.avramukk.com`
and `kubectl -n cloudflare get pods`.
**Escalate to:** yourself; log the incident in [`docs/incident-reports/`](../incident-reports/).

Related: [ADR-008](../adr/008-public-dashboard-cloudflare-tunnel.md),
[ADR-019](../adr/019-cloudflare-tunnel-in-opentofu.md), skill `.opencode/skills/cloudflare-tunnel`.

## Error decoder

| Symptom | Meaning | First check |
|---|---|---|
| `1033` | Tunnel connector not running / not connected | `kubectl -n cloudflare get pods`, `kubectl -n cloudflare logs deploy/cloudflared` |
| `502` | Connector up, no service behind the hostname | tunnel ingress rule + in-cluster Service |
| `530` | Cloudflare can't reach the tunnel origin | same as `1033` |
| DNS resolves but wrong target | stale CNAME | `tofu -chdir=terraform output tunnel_id`, compare with `dig` |

## Verify the tunnel

```bash
export PATH="/opt/homebrew/bin:$PATH"
export KUBECONFIG="$HOME/.kube/config"
kubectl -n cloudflare get pods
kubectl -n cloudflare logs deploy/cloudflared --tail=50 | grep -iE "registered|error"
```

Cloudflare side (connector count, tunnel id):

```bash
cf=terraform/terraform.tfvars   # contains the API token (gitignored)
# list tunnels and their connections via the API or the dashboard (Zero Trust → Networks → Tunnels)
```

## Verify DNS

```bash
dig +short homelab.avramukk.com          # → Cloudflare anycast IPs
# record target must match terraform output tunnel_id
tofu -chdir=terraform output tunnel_id
```

## Restart the connector

```bash
kubectl -n cloudflare rollout restart deploy/cloudflared
kubectl -n cloudflare rollout status deploy/cloudflared
```

## Rotate the tunnel token

1. Re-fetch the token from Cloudflare (API `GET /accounts/{id}/cfd_tunnel/{tid}/token`
   or the dashboard).
2. Re-seal:

```bash
kubectl -n cloudflare get secret cloudflared-tunnel-token -o yaml > /tmp/secret.yaml   # edit token
kubeseal --cert /tmp/ss-cert.pem --format yaml < /tmp/secret.yaml > infra/cloudflared/manifests/sealedsecret.yaml
```
3. Commit and let Argo CD sync; restart the deployment.

## Guardrails

- The tunnel token lives **only** in `SealedSecret`; never commit plaintext.
- Public hostnames are anonymous; never expose admin endpoints (Argo CD, dashboards
  with write access).
- Changing tunnel ingress is `tofu apply` in `terraform/` — review the plan first.
