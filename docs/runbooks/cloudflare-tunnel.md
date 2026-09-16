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

## A hostname resolves publicly but not on one machine

**Symptom:** a public hostname works from other networks but a specific machine
gets `ERR_NAME_NOT_RESOLVED`, while `dig @1.1.1.1 <host>` answers correctly.

**Cause:** local DNS state, not the tunnel. In this lab it was **Tailscale
MagicDNS** — Tailscale had installed its resolver (`100.100.100.100`) as a system
resolver, and that resolver stopped answering (`dig @100.100.100.100 google.com`
returned nothing / `SERVFAIL`). macOS then cached the failed lookup.

**Diagnose, cheapest first:**

```bash
dig +short demo.avramukk.com            # system resolver
dig +short @1.1.1.1 demo.avramukk.com   # public resolver (control)
dig +short @100.100.100.100 google.com  # Tailscale MagicDNS health
scutil --dns | grep -c 100.100.100.100  # is MagicDNS a resolver?
dscacheutil -q host -a name demo.avramukk.com   # is it cached?
```

**Fix (in order):**

```bash
# 1. If MagicDNS is broken, take Tailscale out of the DNS path
tailscale set --accept-dns=false

# 2. Flush the macOS cache (needs sudo)
sudo dscacheutil -flushcache && sudo killall -HUP mDNSResponder

# 3. Clear the browser's own host cache
#    Chrome/Edge: chrome://net-internals/#dns → Clear host cache, then reload
```

The site itself is fine if `dig @1.1.1.1` answers and
`curl --resolve <host>:443:<ip> https://<host>` returns 200.

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
