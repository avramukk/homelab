# Skills

Which agent skills this repository leans on, and when to reach for them.

> Run 1 is a pilot ([ADR-025](adr/025-run-1-pilot.md)) and its documentation is the
> deliverable. This page is the toolkit record so run 2 does not have to
> rediscover which skill applies to which job.
>
> The short, agent-facing version lives in [`../AGENTS.md`](../AGENTS.md); this is
> the longer, human-readable map.

## Use these often

The skills that carry most of the work in this repository.

| Skill | Reach for it when | Where it applies here |
|---|---|---|
| `kubernetes-skill` | writing, reviewing or refactoring any manifest, Helm value or policy | everything under `infra/` |
| `k3d` | cluster lifecycle: create/stop/delete, nodes, image import, kubeconfig | the runtime ([ADR-016](adr/016-runtime-k3s-k3d.md)) |
| `argocd` | sync/diff Applications, debug `OutOfSync`/`Degraded`, app-of-apps | `infra/apps/`, `infra/bootstrap/` |
| `cloudflare-tunnel` | public hostnames, tunnel ingress, DNS, 502/1033 errors | `terraform/`, `infra/cloudflared/` |
| `restic-backup` | repository init, backup, `restic check`, prune, restore drills | `infra/backup/` (pending — issue #2) |
| `terraform-skill` | writing or reviewing OpenTofu for external resources | `terraform/` (Cloudflare, Tailscale) |
| `terraform-plan-review` | **before every** `tofu apply` — risk review of the plan | same |
| `git-workflow-and-versioning` | commits, branches, history, releases, changelog | every change (Conventional Commits) |
| `documentation-and-adrs` | writing an ADR, or deciding what must be documented | `docs/adr/` (001–025) |
| `promql` / `prometheus` | SLO and recording rules, burn-rate math, cardinality hunting | `infra/demo/manifests/slo.yaml`, dashboards |
| `verification-before-completion` | before claiming anything is done, fixed or passing | every task's exit criteria |

## Project skills

These ship **with this repository** (`/.opencode/skills`), so they travel with the
code and take precedence over global skills of the same name.

| Skill | Covers |
|---|---|
| `k3d` | creating, stopping, deleting the cluster; node containers; image imports |
| `argocd` | sync/diff, drift, self-management, recovery of a lost Argo CD |
| `cloudflare-tunnel` | tunnel creation, in-cluster `cloudflared`, DNS, debugging |
| `restic-backup` | Restic repo, backup CronJob, verification, restore procedure |
| `cilium` | NetworkPolicy, default-deny, Hubble, endpoint debugging |
| `talos` | deferred runtime, kept for a future bare-metal host |

## By task

| Task | Skill |
|---|---|
| Any Kubernetes manifest, Helm chart, Kustomize | `kubernetes-skill` |
| NetworkPolicy, CNI, east-west isolation | `cilium`, `kubernetes-security` |
| Cluster is down or misbehaving | `k3d`, `kubernetes-troubleshooting` |
| Deployment is not reconciling | `argocd` |
| Publishing or debugging a public hostname | `cloudflare-tunnel`, `cloudflare` |
| Terraform/OpenTofu change | `terraform-skill`; `terraform-plan-review` before apply |
| Terraform state surgery (mv/import/rm) | `terraform-state-operations` |
| Instrumenting a service (metrics/logs/traces) | `observability-and-instrumentation`, `opentelemetry` |
| Metrics, PromQL, SLO math | `prometheus`, `promql` |
| Logs | `loki` |
| Traces | `tempo` |
| Telemetry collection pipeline | `alloy` |
| Dashboards as code | `dashboarding`, `grafana-oss` |
| Alerts, routing, on-call, SLO burn | `alerting-irm`, `oncall-irm` |
| Probes and load tests | `synthetic-monitoring-checks`, `k6`, `testing` |
| A runbook for 2am | `runbook-creator` |
| An incident | `incident-response`, `systematic-debugging` |
| Security review, hardening, supply chain | `security-and-hardening`, `security-audit` |
| Commits, branches, releases | `git-workflow-and-versioning` |
| ADRs and docs | `documentation-and-adrs` |
| Plan a multi-step change | `writing-plans`, `executing-plans` |
| Backfill or history analysis | `historical-pattern-analysis` |
| Diagrams | `excalidraw-skill` |
| Reading a long web page | `defuddle` |
| Before claiming completion | `verification-before-completion` |
| Ready to ship | `production-readiness`, `shipping-and-launch` |
| CI pipelines | `ci-cd-and-automation` |
| Code review (give or receive) | `code-review-and-quality`, `requesting-code-review`, `receiving-code-review` |

## Observability deep-dives

Use when the observability stack itself is the subject, rather than a change that
merely touches it.

| Skill | For |
|---|---|
| `adaptive-metrics`, `cost-management`, `dpm-finder` | active-series cost, Adaptive Metrics, DPM analysis |
| `prometheus-cardinality-troubleshooter`, `prometheus-label-strategy` | a cardinality fire, or preventing one at the source |
| `loki-label-analyzer` | auditing the Loki label schema |
| `mimir` | long-term metrics storage, multi-tenancy |
| `pyroscope`, `profilecli-insights` | continuous profiling, flame graphs |
| `beyla` | zero-code eBPF instrumentation |
| `fleet-management` | fleet-wide Alloy configuration |
| `infrastructure`, `send-data`, `cloud-integrations` | onboarding hosts/clusters/clouds to Grafana Cloud |
| `assistant-mcp` | wiring an AI agent to a Grafana instance |
| `database-observability` | PostgreSQL/MySQL query and schema visibility |

## Deliberately out of scope

This repository is a single-cluster Kubernetes homelab. Skill families for
application authentication, payments/billing, end-user UI polish, design systems
and personal-assistant workflows are installed on the workstation but are **not
used here** — listing them would add noise, not signal. Reach for them in the
projects they belong to.

## Where skills live

| Location | Scope |
|---|---|
| `.opencode/skills/` | **this repository** — travel with the code, highest precedence |
| `~/.agents/skills/` | the workstation's canonical global store |
| `~/.config/opencode/skills/` | opencode-specific workflow skills |

Invoke a skill by id with the `skill` tool. When a project skill and a global
skill share a name, the project one wins — that is why `k3d`, `argocd` and
`cloudflare-tunnel` here describe *this* cluster rather than the general case.
