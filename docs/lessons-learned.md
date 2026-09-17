# Lessons Learned — Run 1

> Status: `Living`. Run 1 is a **pilot** ([ADR-025](adr/025-run-1-pilot.md)): it
> exists to surface problems and settle decisions before a clean rebuild. This
> register is the input to **run 2**.
>
> Read this **before** the ADRs.

## How to use this document

Every entry has the same shape:

- **What hurt** — the observable symptom.
- **Root cause** — the mechanism, not the person.
- **Run 2** — the concrete change to make.
- **Evidence** — a commit, an ADR, or an issue. Never "we remember".

Add an entry the moment something costs more than an hour, forces a reversal, or
surprises us — while it is still fresh (see [lesson 12](#12-keep-this-register-while-working-not-afterwards)).

## By the numbers (run 1)

| Metric | Value |
|---|---|
| Commits | 80 (2026-09-16 → 2026-09-17) |
| Fix commits | 14 `fix` + 2 `security` |
| Direct reverts | 2 |
| ADRs superseded or amended | 5 of 25 (001, 002, 008, 018, 021) |
| Backfilled issues | 56 — 14 epics, 42 tasks, of which 13 are `bug` |
| Unfinished work carried to run 2 | 6 items |

## The lessons

### 1. Choose the runtime after a spike, not before

**What hurt:** the runtime was decided in planning (Talos) and then blocked twice
on the real host — a QEMU `vmnet` DHCPv4 failure, then a Colima overlayfs
failure — forcing a full replan to k3s via k3d.

**Root cause:** the decision came from a market scan, not from a validation
against the actual macOS + Apple Silicon constraints.

**Run 2:** timebox a runtime spike on the real host *first*; only then write the
ADR.

**Evidence:** [ADR-001](adr/001-runtime-talos-linux.md) → [ADR-016](adr/016-runtime-k3s-k3d.md),
commits `5845c91`, `9e266c7`

### 2. Decide the exposure policy before publishing anything

**What hurt:** an anonymous, read-only Grafana was published and reverted the same
day. The same pattern repeated with the demo API.

**Root cause:** the "showcase value" argument was accepted before the
information-disclosure cost was reviewed.

**Run 2:** write the exposure policy — what may be public, what is tailnet-only —
before the first hostname exists. Any non-read-only or high-information surface is
private by default.

**Evidence:** `46556ac` → `db4fe8b`,
[ADR-008](adr/008-public-dashboard-cloudflare-tunnel.md) → [ADR-020](adr/020-grafana-private.md);
[ADR-021](adr/021-publish-demo-api.md) → [ADR-022](adr/022-demo-api-tailnet-only.md)

### 3. Pick tools that are configurable as code from day one

**What hurt:** Uptime Kuma exposes no admin REST API, so its monitors were
click-ops and the status page had to be migrated to Gatus and rebuilt.

**Root cause:** the tool was chosen for its UI, not for a declarative config path.

**Run 2:** for anything that must be reproducible, require a documented
config-as-code path *before* adoption — otherwise it is a candidate for rejection.

**Evidence:** [ADR-024](adr/024-status-page-gatus.md), issues #51–#56

### 4. Design NetworkPolicy from the real traffic graph

**What hurt:** three separate breakages — demo → PostgreSQL, cloudflared → demo,
Gatus → demo/Grafana.

**Root cause:** policies were written per-service without a whole-cluster traffic
picture. k3s **enforces** NetworkPolicy (kube-router), so a wrong rule is an
outage, not a no-op.

**Run 2:** draw the traffic graph first (namespace → namespace, port, direction),
generate the policies from it, and smoke-test every path after applying.

**Evidence:** issues #41, #42, #54, #55

### 5. Build multi-arch images from the first CI run

**What hurt:** the demo image was built single-arch and would not start on the
arm64 cluster — `no match for platform in manifest`.

**Root cause:** CI runners are amd64 while the cluster is arm64.

**Run 2:** `linux/amd64,linux/arm64` in the very first `build-push-action`, with
buildx + QEMU configured from the start.

**Evidence:** issue #39, commit `e634b1e`

### 6. One owner per namespace

**What hurt:** several Argo CD Applications declared the same Namespace, so it
flapped between owners and one app reported duplicates.

**Root cause:** `CreateNamespace=true` in sync options *plus* a separate bare
`Namespace` manifest.

**Run 2:** a single `namespaces` Application owns every namespace; child
Applications never create them.

**Evidence:** issue #22, commits `69d4840`, `436c97d`

### 7. Validate Helm values against the chart schema

**What hurt:** Traefik rejected a values key (`logs`), and Loki's ServiceMonitor
sat outside the chart's `monitoring` block, so it was never scraped.

**Root cause:** values were written from memory and examples rather than validated
against the chart.

**Run 2:** render and schema-validate every chart in CI (`helm template`,
`kubeconform`) before merge — not only for the resources we author by hand.

**Evidence:** issues #25, #36

### 8. Do not assume IP allowlists work behind the k3d load balancer

**What hurt:** a Traefik `ipAllowList` restricted to `100.64.0.0/10` blocked
legitimate traffic and had to be reverted.

**Root cause:** the k3d load-balancer proxy **NATs the client address**, so
Traefik never sees the tailnet IP.

**Run 2:** establish what the ingress actually observes (`remote_addr`,
`X-Forwarded-For`) before writing any address-based rule; prefer restricting at
the network layer over header/address rules.

**Evidence:** issue #59, [ADR-023](adr/023-host-ports-lan-reachable.md),
commits `e5b8c47` → `fce0960`

### 9. Ship alert delivery together with the alert rules

**What hurt:** SLO recording and multi-window burn-rate rules exist, but
Alertmanager still has a **null receiver** — if the service burned its error
budget today, nobody would be paged.

**Root cause:** the rules and the delivery path were split into separate work
items, and only the "interesting" half shipped.

**Run 2:** one definition of done for alerting — rule **+** receiver **+** routing
**+** a test fire. An alert that has never fired is not done.

**Evidence:** [ADR-009](adr/009-alerting-telegram-slo.md), forward issue #1

### 10. Backups are part of "production-grade", not a follow-up

**What hurt:** the system was described as production-grade while the backup
design existed only on paper and no restore had ever been run.

**Root cause:** the visible work (GitOps, observability) crowded out durability.

**Run 2:** no workload is "done" until its data has a tested restore path; the
restore drill is scheduled together with the feature that creates the data.

**Evidence:** [ADR-010](adr/010-backup-restic-local.md), forward issues #2 and #5

### 11. Put process guardrails in before the first commit

**What hurt:** all 80 commits went straight to `main` — no pull requests, no
branch protection, no tags, no releases.

**Root cause:** the repository started as a personal scratchpad and only later
became a public showcase, so the process was never imposed.

**Run 2:** branch protection, required PRs and required status checks from commit
one; tag the first phase; keep history linear and signed from the start.

**Evidence:** forward issues #62, #4

### 12. Keep this register while working, not afterwards

**What hurt:** this document had to be **reconstructed from `git log`** — 80
commits had to be re-read to recover why things happened.

**Root cause:** there was no habit of writing the lesson at the moment of
friction.

**Run 2:** add an entry when something costs more than an hour, forces a reversal,
or surprises us — before moving to the next task.

**Evidence:** this file

## ADR reversals (the record)

| Original | Superseded / amended by | Why |
|---|---|---|
| [ADR-001](adr/001-runtime-talos-linux.md) — Talos Linux runtime | [ADR-016](adr/016-runtime-k3s-k3d.md) — k3s via k3d | blocked by macOS/Colima twice |
| [ADR-002](adr/002-cluster-topology.md) — provisioning model | [ADR-016](adr/016-runtime-k3s-k3d.md) | runtime pivot; node count unchanged |
| [ADR-008](adr/008-public-dashboard-cloudflare-tunnel.md) — public dashboard incl. Grafana | [ADR-020](adr/020-grafana-private.md) — Grafana private | information disclosure |
| [ADR-018](adr/018-cloudflare-tunnel-cli-opentofu-dns.md) — tunnel via CLI | [ADR-019](adr/019-cloudflare-tunnel-in-opentofu.md) — tunnel in OpenTofu | keep external state in one tool |
| [ADR-021](adr/021-publish-demo-api.md) — public demo API | [ADR-022](adr/022-demo-api-tailnet-only.md) — tailnet-only | open write endpoint |
| [ADR-008](adr/008-public-dashboard-cloudflare-tunnel.md) — Uptime Kuma status page | [ADR-024](adr/024-status-page-gatus.md) — Gatus | no admin API → click-ops |

> [ADR-023](adr/023-host-ports-lan-reachable.md) supersedes nothing — it
> *qualifies* ADR-020 by recording an accepted risk.

## Unfinished in run 1 → run 2 scope

| Item | Where |
|---|---|
| Alert delivery to Telegram | [ADR-009](adr/009-alerting-telegram-slo.md) · issue #1 |
| Backups and a tested restore | [ADR-010](adr/010-backup-restic-local.md) · issue #2 |
| Trivy, Kyverno, NetworkPolicy/PSS rollout | [ADR-012](adr/012-security-baseline.md) · issue #3 |
| Docs truth-sync, release and first tag | issue #4 |
| Chaos drill and postmortem | [ADR-014](adr/014-incident-drill.md) · issue #5 |
| Branch protection and required PRs | issue #62 |
