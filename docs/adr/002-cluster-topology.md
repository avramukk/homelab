# ADR-002: Cluster topology — 1 control-plane + 2 workers

## Status
Accepted — provisioning superseded by [ADR-016](016-runtime-k3s-k3d.md); the node
count and roles below are unchanged.

## Date
2026-09-16

## Context

With Talos chosen ([ADR-001](001-runtime-talos-linux.md)), the next decision is
how many nodes. The host has 14 vCPU and 24 GiB RAM total, shared with macOS and
the QEMU process overhead.

Constraints:

- Real multi-node behaviour (scheduling, node roles, `nodeSelector`, taints) is
  part of the point — a single node does not demonstrate it.
- Control-plane redundancy requires an odd number of etcd members (3 minimum).
- RAM is the binding constraint; each node carries a full OS + kubelet + runtime.

## Decision

Run **one control-plane node and two workers**, each sized at approximately
2 vCPU / 4 GiB.

## Alternatives considered

### Single node
- Pros: least RAM, simplest.
- Cons: no scheduling story, no node roles, no multi-node failure modes; defeats
  the purpose of choosing Talos.
- Rejected.

### 3 control-plane + 2 workers (HA etcd)
- Pros: real etcd quorum, control-plane HA, strongest topologies.
- Cons: five VMs at 3–4 GiB each ≈ 15–20 GiB, leaving macOS very little headroom;
  QEMU overhead on top; likely swapping, which makes the demo slow and flaky.
- Rejected for now: **the upgrade path**, not the ceiling. A third
  control-plane is added when hardware allows.

### 1 control-plane + 1 worker
- Pros: more headroom per node.
- Cons: one worker is not enough to show scheduling/spread behaviour.
- Rejected.

## Consequences

- **Accepted limitation:** a single control-plane means a **single etcd member** —
  no control-plane HA. This is documented explicitly rather than hidden.
- **Upgrade path:** move to 3 control-planes if the host gains RAM, with no change
  to anything above the cluster layer.
- **Capacity:** three VMs consume roughly 12–15 GiB, leaving headroom for macOS
  under normal load.
- Cluster-level concerns (storage, network policy) are handled elsewhere:
  [ADR-010](010-backup-restic-local.md), [ADR-012](012-security-baseline.md).
