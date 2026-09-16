# Credentials

Operator credentials for this homelab live in **1Password**, not in the
repository. Git holds only `SealedSecret` ciphertext ([ADR-005](../adr/005-secrets-sealed-secrets.md)).

## Where things live

| Store | Holds |
|---|---|
| **1Password** (Personal vault, items prefixed `Homelab ·`) | human-facing credentials: admin logins, DB password, Cloudflare tokens |
| **Git** | `SealedSecret` ciphertext only — decryptable solely by the in-cluster controller |
| **Cluster** | the decrypted `Secret`s, mounted into workloads |
| **Host** | `~/.cloudflared/` credentials, `~/.kube/config`, `~/.talos/` — never committed |

## Inventory (1Password → Personal)

| Item | Contents | Used by |
|---|---|---|
| `Homelab · Overview` | public URLs, repo link, cluster summary | you |
| `Homelab · Grafana` | `admin` login + port-forward command | Grafana (private, Tailscale) |
| `Homelab · Argo CD` | `admin` login + port-forward command | Argo CD (private, Tailscale) |
| `Homelab · PostgreSQL (demo)` | `demo` DB user + in-cluster address | demo service / manual psql |
| `Homelab · Cloudflare API token` | API token for `avramukk.com` | OpenTofu (`terraform/`) |
| `Homelab · Cloudflare Tunnel token` | tunnel `homelab` token | `cloudflared` (`SealedSecret`) |

Find them fast:

```bash
op item list --vault Personal | grep Homelab
op item get "Homelab · Grafana" --vault Personal
```

## Rotation

| Credential | When / how |
|---|---|
| Cloudflare API token | short-lived by design; recreate in the dashboard, update `terraform/terraform.tfvars`, `tofu apply`; update the 1Password item |
| Cloudflare Tunnel token | re-fetch, re-seal: `kubeseal --cert <cert> --format yaml < secret.yaml > infra/cloudflared/manifests/sealedsecret.yaml`, commit, restart the deployment |
| Grafana admin | change in Grafana UI, update the `grafana-admin` SealedSecret, update 1Password |
| Argo CD admin | change in the UI, then delete `argocd-initial-admin-secret` |
| PostgreSQL | `ALTER USER demo PASSWORD`, re-seal `postgres-credentials`, update 1Password |

## Re-sealing a secret (no plaintext in Git)

```bash
kubeseal --fetch-cert --controller-name sealed-secrets-controller \
  --controller-namespace kube-system > /tmp/ss-cert.pem

kubectl -n <ns> create secret generic <name> \
  --from-literal=key=value --dry-run=client -o yaml \
  | kubeseal --cert /tmp/ss-cert.pem --format yaml > infra/<path>/<name>-sealedsecret.yaml
```

## Rules

- **Never** commit plaintext credentials; gitleaks blocks it (pre-commit + CI).
- Prefer short-lived tokens; prefer in-cluster credentials over host files.
- The Sealed Secrets **sealing key** is the master secret: back it up encrypted
  out of band ([ADR-010](../adr/010-backup-restic-local.md)).
