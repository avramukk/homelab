# Homelab

A production-grade SRE homelab, built in public. Single-cluster Kubernetes,
operated with GitOps, observability, SLOs, and documented decisions — **run 1, a
pilot**.

> **Run 1 — a pilot.** This is the first, deliberately experimental build. It
> exists to surface problems and settle decisions before a clean rebuild (run 2).
> Every mistake, reversal and fix is documented — start with
> [`docs/lessons-learned.md`](docs/lessons-learned.md) and
> [ADR-025](docs/adr/025-run-1-pilot.md).

> Status: **Running** — a single-node k3s cluster, operated with production-grade
> practice but explicitly **not** production. See
> [`docs/architecture/overview.md`](docs/architecture/overview.md).

## Start here

| I want to… | Go to |
|---|---|
| Understand the system | [`docs/architecture/overview.md`](docs/architecture/overview.md) |
| See **why** things are built this way | [`docs/adr/`](docs/adr/) |
| Operate it (SLOs, backups, on-call) | [`docs/operations/`](docs/operations/) |
| Fix something at 2am | [`docs/runbooks/`](docs/runbooks/) |

## How it is built

- **Runtime:** k3s via k3d, 1 server + 2 agents ([ADR-016](docs/adr/016-runtime-k3s-k3d.md), [ADR-002](docs/adr/002-cluster-topology.md))
- **GitOps:** Argo CD ([ADR-003](docs/adr/003-gitops-argocd.md))
- **CI:** GitHub Actions ([ADR-004](docs/adr/004-ci-github-actions.md))
- **Observability:** self-hosted LGTM + OpenTelemetry ([ADR-007](docs/adr/007-observability-lgtm-otel.md))
- **Public edge:** Cloudflare Tunnel ([ADR-008](docs/adr/008-public-dashboard-cloudflare-tunnel.md))
- **Backups:** Restic, tested with a monthly restore drill ([ADR-010](docs/adr/010-backup-restic-local.md))

## Conventions

- History is **linear** and follows Conventional Commits; one ADR per commit.
- **No secrets in Git** — Sealed Secrets ciphertext and signed commits only.
- See [`CONTRIBUTING.md`](CONTRIBUTING.md) and [`docs/README.md`](docs/README.md).
