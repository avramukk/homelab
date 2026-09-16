# ADR-001: Cluster runtime — Talos Linux over Colima/k3s

## Status
Accepted

## Date
2026-09-16

## Context

The homelab runs on a single Apple Silicon laptop (`lab-host`, M4, 24 GiB) whose
host OS is macOS. macOS cannot run Kubernetes natively, so a virtualization layer
is required. This is the foundation decision: everything above it (GitOps,
observability, backups) is unaffected by the choice, but the *credibility* of the
showcase — this is an SRE portfolio artifact, not a toy — depends on it.

Constraints:

- One physical host, no additional hardware budget.
- Must look and behave like production: real nodes, real upgrade story, real
  failure modes.
- Operator is Solo; operational complexity must stay runnable.

## Decision

Use **Talos Linux**, running as QEMU virtual machines provisioned with
`talosctl cluster create`. Talos is an immutable, API-driven Kubernetes OS: no
shell, no SSH, atomic A/B image upgrades with automatic rollback.

A Talos deployment was already installed and evaluated: `colima` + `docker` +
`lima` on `lab-host`. Colima is **superseded** and left installed but unused.

## Alternatives considered

### k3s on Colima (single node)
- Pros: already installed; fastest path; bundles Traefik/CoreDNS/local-path.
- Cons: single node with SQLite/kine, no etcd, no scheduling/HA patterns; a
  reviewer sees "one container runtime plus k3s in a VM" — reads as a toy.
- Rejected: weakest credibility for an SRE showcase.

### k3d (k3s in Docker), multi-node
- Pros: multiple nodes on one host cheaply; HA topologies possible.
- Cons: nodes are containers, not machines; requires a second virtualization
  layer (Colima→dockerd→k3d); cluster does not survive a host reboot cleanly.
- Rejected: better than single-node, but still not "real nodes".

### Ubuntu VMs + kubeadm
- Pros: the classic on-prem install path (`kubeadm join`, certs, etcd) — maximum
  "enterprise" optics.
- Cons: heaviest to build and maintain; slowest feedback loop; the OS layer itself
  teaches nothing new.
- Rejected: cost/benefit worse than Talos for the same credibility.

### Docker Compose / Podman quadlets
- Pros: honest and simple; most homelabs need no Kubernetes.
- Cons: no Kubernetes surface at all — no platform skills, no CKA/CKS alignment.
- Rejected: the workload catalog is explicitly Kubernetes-centric.

## Consequences

- **Gained:** real nodes, immutable OS, declarative machine config, atomic
  upgrades, and strong alignment with current "production-correct homelab"
  practice.
- **Accepted:** nodes have no shell — all node-level operations go through the
  Talos API (`talosctl`); this is a skill, but it is a different operating model.
- **Security burden:** Talos machine secrets and `talosconfig` are cluster
  identity material and must never enter the public repository (see
  [ADR-005](005-secrets-sealed-secrets.md)).
- **Cost:** QEMU on Apple Silicon has overhead; node sizing is bounded by
  24 GiB (see [ADR-002](002-cluster-topology.md)).
- **Colima** remains installed on the host, unused.
