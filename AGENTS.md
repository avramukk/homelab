# AGENTS.md — homelab

Repo-specific instructions for agents. Everything below was verified against the
repository; prefer the executable source (config, workflows, manifests) over prose.

## What this repo is

A **single-node k3s cluster** (k3d on Colima, on a Mac host) managed entirely
from Git by Argo CD, plus the documentation around it. It is a public SRE
showcase: keep it clean, English-only, and free of secrets or employer names.

```
Internet ──HTTPS──► Cloudflare (DNS + Tunnel)          GitHub ──pull──► Argo CD
                        ▲                                    ▲
                        │ outbound-only                      │
        ┌───────────────┴──────────────  lab-host (Apple Silicon)  ──────────┐
        │  Colima VM (Linux) → k3s (k3d): 1 server + 2 agents                │
        │    cloudflared · Traefik · demo+postgres · landing · uptime-kuma   │
        │    Prometheus · Loki · Tempo · Grafana · Alloy · Argo CD           │
        └────────────────────────────────────────────────────────────────────┘
                        ▲
        Tailscale (private) ── ssh/port-forward ── developer laptop
```

## Read before acting

- `docs/adr/` — every significant decision. **Append-only**: to change a decision,
  add a new ADR that supersedes the old one; never rewrite history of an ADR.
- `docs/architecture/` — how the system is wired (mermaid diagrams live here).
- `infra/apps/*.yaml` — the Argo CD Applications; the root app-of-apps points at
  `infra/apps` (`infra/bootstrap/root-app.yaml`).
- `.github/workflows/` — CI: `gitleaks.yml` (secrets) and `build-demo-app.yml`.

## Layout

| Path | Holds |
|---|---|
| `apps/demo/` | Go service (module `github.com/avramukk/homelab/demo`) |
| `infra/apps/` | one Argo CD `Application` per component (app-of-apps) |
| `infra/<component>/` | Helm values and manifests for that component |
| `terraform/` | OpenTofu for **external** resources: Cloudflare DNS + Tunnel |
| `docs/` | ADRs, architecture, runbooks, operations, security |
| `.opencode/skills/` | project skills (see “Skills” below) |

## Commands (exact)

The cluster runs on the **`lab-host`** machine, not on the developer laptop:

```bash
ssh lab-host
export PATH="/opt/homebrew/bin:$PATH"
export DOCKER_HOST="unix://$HOME/.colima/default/docker.sock"
kubectl config use-context k3d-homelab
kubectl get nodes
```

Force Argo CD to reconcile now (after pushing):

```bash
kubectl -n argocd annotate application root argocd.argoproj.io/refresh=hard --overwrite
```

Demo service checks (on the laptop):

```bash
cd apps/demo && go vet ./... && go test ./...
```

Secret scanning (must be clean before pushing):

```bash
git config core.hooksPath .githooks   # once per clone — enables the pre-commit hook
gitleaks git .
```

OpenTofu (token lives in `terraform/terraform.tfvars`, gitignored):

```bash
cd terraform && tofu init && tofu plan && tofu apply
```

Re-seal a secret (never commit plaintext):

```bash
kubeseal --fetch-cert --controller-name sealed-secrets-controller --controller-namespace kube-system > /tmp/ss-cert.pem
kubectl -n <ns> create secret generic <name> --from-literal=K=V --dry-run=client -o yaml \
  | kubeseal --cert /tmp/ss-cert.pem --format yaml > infra/<path>/<name>-sealedsecret.yaml
```

## Hard-won gotchas

1. **No push-deploy.** CI only validates and builds images. Argo CD pulls from Git;
   change Git, not the cluster. `kubectl apply` / `helm install` are bootstrap-only.
2. **Images must be multi-arch.** The cluster is arm64 (Apple Silicon) while CI runners
   are amd64 → `build-demo-app.yml` builds `linux/amd64,linux/arm64`. A single-arch image
   fails with `no match for platform in manifest`.
3. **k3s enforces NetworkPolicy** (its kube-router controller), unlike plain flannel.
   A namespace `default-deny` blocks new pods until you add an allow rule — e.g. ingress
   from `cloudflare` (cloudflared), `observability` (Prometheus), `traefik`.
4. **Storage is `local-path` only** — no replication; durability comes from backups (ADR-010).
5. **Two ingress paths, deliberately separate:**
   - *Public* → Cloudflare Tunnel → `cloudflared` → Service. Only `homelab` (landing) and
     `status` (Uptime Kuma) are public today.
   - *Private* → DNS-only `A` record pointing at the **tailnet IP**, routed by Traefik on
     the tailnet with Let's Encrypt (DNS-01) certs. `grafana`, `argocd`, `demo` live here.
   Do not expose admin UIs (Grafana, Argo CD) publicly — ADR-020.
6. **Hostnames are managed in `terraform/`** via `public_hostnames`, `private_hostnames`,
   `reserved_hostnames`. Do not add Cloudflare DNS records by hand.
7. **Never commit** `terraform.tfvars`, `*.tfstate`, `INBOX.md`, `TODO-*.md`, keys, tokens,
   or build artifacts (a 28 MB Go binary was once committed here — see `.gitignore`).
   No employer or client names anywhere, including history.
8. **History is linear and signed.** Commits are SSH-signed (repo-local config); one ADR per
   commit; Conventional Commits. Rewriting history (filter-repo) requires re-signing every
   commit and force-pushing.
9. **Tailscale DNS can wedge on the laptop**: MagicDNS stops answering, so `dig` works but
   curl/browsers fail with `ERR_NAME_NOT_RESOLVED`. Fix: `tailscale set --accept-dns=false`
   then flush (see `docs/runbooks/cloudflare-tunnel.md`).

## Skills

Project skills in `.opencode/skills/` — load them for the matching tool:
`k3d` (cluster lifecycle), `argocd` (sync/debug/self-management),
`cloudflare-tunnel` (public hostnames, DNS), `restic-backup` (backups/restore),
`talos` (deferred runtime — kept for a future bare-metal host).

The full human-readable map of which skills this repository uses, and when, is in
[`docs/skills.md`](docs/skills.md).

Global skills worth loading by task:

| Task | Skill |
|---|---|
| Writing/reviewing any Kubernetes manifest, Helm, Kustomize | `kubernetes-skill` (KubeShark failure-mode workflow) |
| Terraform/OpenTofu | `terraform-skill`; `terraform-plan-review` before every apply |
| Commits, branches, releases, history | `git-workflow-and-versioning` |
| ADRs and docs | `documentation-and-adrs` |
| Instrumentation (metrics/logs/traces) | `observability-and-instrumentation` |
| Runbooks, incident response | `runbook-creator`, `incident-response`, `alerting-irm` |
| Security review, hardening | `security-and-hardening`, `security-audit` |
| Cloudflare platform / Zero Trust | `cloudflare`, `cloudflare-one` |
| NetworkPolicy / CNI | `cilium` |
| Metrics, logs, traces, collection | `prometheus`, `promql`, `loki`, `tempo`, `alloy` |
| Before claiming completion | `verification-before-completion` |
| Debugging an unknown failure | `systematic-debugging` |

## Conventions

- **Conventional Commits**, one logical change per commit; an ADR is one commit.
- **Docs change with the system**, in the same PR.
- **PRs** state Decision / Alternatives / Validation / Rollback (`.github/pull_request_template.md`).
- Run `gitleaks` before every push; CI enforces it on push/PR and weekly.
