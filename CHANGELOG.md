# Changelog

All notable changes to this homelab. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and releases are tagged
by phase.

## [Unreleased]

### Added
- Documentation architecture (`docs/`).
- ADRs 001–016.
- Phase 1: multi-node Kubernetes cluster — k3s v1.35.5 via k3d (1 server + 2 agents) on Colima.
- Phase 2: Argo CD bootstrapped (chart `argo-cd` 10.9.1) and self-managed from Git; app-of-apps root in `infra/apps` ([ADR-017](docs/adr/017-argocd-bootstrap.md)).
- Phase 3 (partial): Traefik ingress (chart 41.6.0) and Sealed Secrets controller (v0.40.0), both GitOps-managed.
- External: Cloudflare Tunnel `homelab` and `homelab.avramukk.com` fully managed by OpenTofu ([ADR-019](docs/adr/019-cloudflare-tunnel-in-opentofu.md), supersedes [ADR-018](docs/adr/018-cloudflare-tunnel-cli-opentofu-dns.md)).
- cloudflared runs in-cluster (2 connectors, token in a `SealedSecret`); `homelab.avramukk.com` is live (catch-all 404 until a service is published).
- Phase 3: observability stack — kube-prometheus-stack, Loki, Tempo, Alloy (logs→Loki, OTLP→Tempo), Grafana datasources as code, Uptime Kuma.
- Published `homelab.avramukk.com` (Grafana, anonymous read-only) and `status.avramukk.com` (Uptime Kuma) through the tunnel.

### Changed
- Runtime pivoted from Talos-on-QEMU to k3s-via-k3d after macOS networking and
  kernel incompatibilities ([ADR-016](docs/adr/016-runtime-k3s-k3d.md) supersedes
  [ADR-001](docs/adr/001-runtime-talos-linux.md)).
- Architecture Decision Records 001–015.
