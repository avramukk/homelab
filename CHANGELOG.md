# Changelog

All notable changes to this homelab. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and releases are tagged
by phase.

## [Unreleased]

### Added
- Documentation architecture (`docs/`).
- ADRs 001–016.
- Phase 1: multi-node Kubernetes cluster — k3s v1.35.5 via k3d (1 server + 2 agents) on Colima.

### Changed
- Runtime pivoted from Talos-on-QEMU to k3s-via-k3d after macOS networking and
  kernel incompatibilities ([ADR-016](docs/adr/016-runtime-k3s-k3d.md) supersedes
  [ADR-001](docs/adr/001-runtime-talos-linux.md)).
- Architecture Decision Records 001–015.
