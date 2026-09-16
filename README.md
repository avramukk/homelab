# Homelab

A production-grade SRE homelab, built in public. Single-cluster Kubernetes on
immutable infrastructure, operated with GitOps, observability, SLOs, tested
backups, and documented decisions.

> Status: **planning → scaffolding.** No infrastructure is running yet. Decisions
> are captured as ADRs; implementation follows phase by phase.

## Start here

| I want to… | Go to |
|---|---|
| Understand the system | [`docs/architecture/overview.md`](docs/architecture/overview.md) |
| See **why** things are built this way | [`docs/adr/`](docs/adr/) |
| Operate it (SLOs, backups, on-call) | [`docs/operations/`](docs/operations/) |
| Fix something at 2am | [`docs/runbooks/`](docs/runbooks/) |

## How it is built

- **Runtime:** Talos Linux, 1 control-plane + 2 workers ([ADR-001](docs/adr/001-runtime-talos-linux.md), [ADR-002](docs/adr/002-cluster-topology.md))
- **GitOps:** Argo CD ([ADR-003](docs/adr/003-gitops-argocd.md))
- **CI:** GitHub Actions ([ADR-004](docs/adr/004-ci-github-actions.md))
- **Observability:** self-hosted LGTM + OpenTelemetry ([ADR-007](docs/adr/007-observability-lgtm-otel.md))
- **Public edge:** Cloudflare Tunnel ([ADR-008](docs/adr/008-public-dashboard-cloudflare-tunnel.md))
- **Backups:** Restic, tested with a monthly restore drill ([ADR-010](docs/adr/010-backup-restic-local.md))

## Conventions

- History is **linear** and follows Conventional Commits; one ADR per commit.
- **No secrets in Git** — Sealed Secrets ciphertext and signed commits only.
- See [`CONTRIBUTING.md`](CONTRIBUTING.md) and [`docs/README.md`](docs/README.md).
