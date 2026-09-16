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

## Colima · `k8s` · `researching`
- **Added:** 2026-09-16
- **What it is:** Lima-based container/VM runtime for macOS; supports a docker context and **K3s inside** (`colima start --kubernetes`).
- **Why it matters:** candidate virtualization layer for the cluster on `lab-host` (macOS/arm64).
- **Sources:** https://github.com/abiosoft/colima · https://abiosoft.github.io/colima/
- **Open questions:** vCPU/RAM to give the VM; K3s mode vs Docker-only; stability of 28-day uptime running K8s on an M4.
- **State:** installed `colima 0.10.3` + `docker CLI 29.8.1` + `lima 2.2.0` on `lab-host` (brew, `/opt/homebrew`). **VM not started** — launch deferred by decision.

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