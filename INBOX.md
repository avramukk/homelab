# INBOX — raw Home Lab backlog

This is the decision backlog. Technologies, links, and ideas dropped into the
chat **without context** land here, so nothing gets lost and everything can
later feed the architecture.

## Rules

- One entry = one block below. Newest on top.
- Status is separate from the fact of being recorded.
- `adopted` entries become input to the final plan / architecture.

## Statuses

| Status | Meaning |
|---|---|
| `inbox` | recorded, not yet reviewed |
| `researching` | actively being studied, open questions |
| `adopted` | decided to use |
| `rejected` | reviewed and dropped (with reason) |

## Drops

## Talos Linux · `k8s` · `adopted`
- **Added:** 2026-09-16
- **What it is:** Immutable Kubernetes OS (Sidero Labs): no shell/SSH, atomic A/B upgrades, upstream K8s, managed via `talosctl` API. Ran via `talosctl cluster create` (QEMU, CLI, scriptable) on `lab-host`.
- **Why it matters:** chosen for "production-grade, not a toy" SRE showcase; real VM nodes, production-correct OS story.
- **Sources:** https://www.talos.dev/ · https://www.siderolabs.com/blog/talos-linux-vs-k3s
- **Open questions:** Talos `talosconfig`/machine secrets must stay out of the public repo (SOPS/local). Manual upgrades + runbook chosen over system-upgrade-controller.
- **Decision:** 1 control-plane + 2 workers, ~4 GiB/node, VMs on `lab-host`; Colima superseded (kept installed, unused).

## Target architecture (grill 2026-09-16) · `planning` · `adopted`
- **Added:** 2026-09-16
- **What it is:** Confirmed SRE-showcase plan: single-cluster Talos (1 CP + 2 workers) → Argo CD (GitOps) + GitHub Actions (validate/docs/Renovate) + Sealed Secrets + self-hosted LGTM (Prometheus/Loki/Grafana/Tempo + OTel) + Traefik ingress + Uptime Kuma + Cloudflare Tunnel (public Grafana+status, anonymous, own domain) + custom Go CRUD/Postgres demo (OTel) + Telegram alerts (page/ticket) + SLO/burn-rate + in-cluster Restic→local disk (RPO 24h/RTO ≤4h, monthly restore drill; B2 = known limitation) + OpenTofu (Cloudflare DNS, Tailscale, B2 later) + security hardening (PSS restricted, non-root, Cilium/NetworkPolicy, Trivy in CI, Kyverno) + chaos/postmortem + monorepo (apps/infra/terraform/docs: ADRs+runbooks+mermaid).
- **Sources:** this planning session (grill rounds 1–6).
- **Open questions:** pending full docs skeleton + ADR-001 (runtime decision); domain name to be provided (exists, used for tunnel hostnames).

## Market scan 2026 — popular skills/tools among engineers · `reference` · `researching`
- **Added:** 2026-09-16
- **What it is:** Evidence-based landscape scan (CNCF 2024/2025 surveys, Argo CD survey, Grafana Observability Survey 2025, Stack Overflow 2025, GitHub/community). Full sourced report available in-session.
- **Key signals for us:**
  - GitOps mainstream (77% orgs); **Argo CD** leader (~60% of K8s clusters) vs Flux 17%.
  - ⚠️ **ingress-nginx RETIRED Mar 2026** → Gateway API / **Traefik** is the migration path (drop-in compat layer exists).
  - **k3s** = homelab default (bundles Traefik/ServiceLB/local-path); **Talos Linux** = immutable "production-correct" upgrade path (A/B updates); k0s = clean middle.
  - Observability: **LGTM** (Loki/Grafana/Tempo/Mimir) + **OpenTelemetry** rising (41% prod) + **Alloy** collector; **VictoriaMetrics** = lightweight LTS favourite.
  - Image/dep updates: **Renovate** (native ArgoCD/Flux/Helm managers) > Dependabot.
  - Policy: **Kyverno** graduated (Mar 2026), CEL-based, overtook Gatekeeper.
  - Secrets: external-secrets / sealed-secrets / SOPS (tbd).
  - Backup: **Restic** (S3/B2) + **Borg**; 3-2-1: NAS + offsite object storage; **Longhorn** = K8s storage sweet spot (Ceph overkill for homelab).
  - IaC: **OpenTofu** (MPL-2.0 fork, needs Terraform lacks) vs Terraform (BSL).
  - Skills: **CKA/CKS/CKAD** benchmark; CKS now reinstates CKA (Jun 2026); OpenTelemetry fastest-growing; GenAI on K8s 66% orgs; **Ollama** = local LLM entry.
  - Popular self-hosted apps: Home Assistant (2M+ installs), Immich, Vaultwarden, Jellyfin+*arr, Authentik, Paperless-ngx, Uptime Kuma, AdGuardHome, Forgejo/Gitea+Woodpecker.
- **Why it matters:** feeds stack choices during planning; several picks above (Traefik, GitOps, Renovate, restic) align with Daryl Lundy blueprint + Mischa course.
- **Sources:** https://www.cncf.io/reports/cncf-annual-survey-2024/ · https://www.cncf.io/announcements/2025/07/24/cncf-end-user-survey-finds-argo-cd-as-majority-adopted-gitops-solution-for-kubernetes/ · https://kubernetes.io/blog/2025/11/11/ingress-nginx-retirement · https://grafana.com/observability-survey/2025/ · https://survey.stackoverflow.co/2025/technology
- **Open questions:** which items we adopt (decision during planning) — this entry is the raw input.

## Kubernetes Homelab Course (Mischa — KubeCraft) · `learning` · `researching`
- **Added:** 2026-09-16
- **What it is:** Paid Skool course by Mischa, "The Long Awaited Kubernetes Homelab Course": k8s distribution selection, GitOps deployments, container security, securely exposing apps, automated image updates, Traefik ingress. 15 modules / 3+ hours released, more to come.
- **Why it matters:** Maps directly onto our planned stack (GitOps + Traefik + secure exposure) and the CKA/CKS learning path. Full course materials (community backup, transcripts, homelab repo: Talos Linux + Flux + staging/prod clusters) already live locally at `~/repos/mischa`.
- **Sources:** https://www.skool.com/kubecraft/classroom/6463a665 · local: `~/repos/mischa`
- **Open questions:** 62-day membership gate vs annual; does his Flux + Talos model fit our macOS/Colima target, or is k3s enough?

## SRE homelab blueprint (Daryl Lundy) · `reference` · `inbox`
- **Added:** 2026-09-16
- **What it is:** Blog post — a production-grade SRE homelab as a portfolio artifact on Apple Silicon: Traefik + GitLab + Gitea + shared PostgreSQL + MinIO + Prometheus/Grafana/Loki/Alertmanager, 15+ containers via Docker Compose with explicit limits (~14.5 GB reserved on 32 GB), Restic backups to Synology NAS, mkcert TLS for `*.dev.local`.
- **Why it matters:** The closest worked example to our context (Apple Silicon, self-hosted CI + observability + backups); concrete RAM budgets and documented trade-offs worth stealing. k3s migration is his stated next step.
- **Sources:** https://www.daryllundy.com/blog/building-a-production-grade-sre-homelab-not-just-it-works-on-my-machine
- **Open questions:** GitLab vs lighter CI (Gitea/Drone); single-node Compose vs k3s from day one; offsite backup target (B2).

## Colima · `k8s` · `rejected`
- **Added:** 2026-09-16
- **What it is:** Lima-based container/VM runtime for macOS; supports a docker context and **K3s inside** (`colima start --kubernetes`).
- **Why rejected:** superseded by **Talos Linux** (decided in planning) — single-node k3s on Colima reads as a toy next to a 3-node immutable-OS cluster; Talos gives real nodes + production-correct story.
- **Sources:** https://github.com/abiosoft/colima · https://abiosoft.github.io/colima/
- **State:** installed `colima 0.10.3` + `docker CLI 29.8.1` + `lima 2.2.0` on `lab-host` (brew). **Kept installed, not used** (decision: leave in place).

## Cluster host · `infra` · `inbox`
- **Added:** 2026-09-16
- **What it is:** the cluster target machine — laptop over SSH `lab-host` = `operator@lab-host` (Tailscale), macOS Darwin 25.6 **arm64 (M4)**, hostname `lab-host`, 28 days uptime.
- **Why it matters:** the platform is decided before the tech stack — K8s on macOS requires a virtualization layer.
- **Sources:** `~/.ssh/config` → `Host lab-host`
- **Open questions:** dedicated K8s node or also a general host? Which layer: k3d / kind / OrbStack / VM? External access via Tailscale?

<!-- Add new entries above this line. Copy the template below. -->

---

## Template

```md
## <Name> · `category` · `inbox`
- **Added:** YYYY-MM-DD
- **What it is:** 1–2 sentences
- **Why it matters:** …
- **Sources:** [docs](…) · [github](…)
- **Open questions:** …
```