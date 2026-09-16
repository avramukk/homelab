# ADR-016: Cluster runtime — k3s via k3d (Talos deferred)

## Status
Accepted

## Date
2026-09-16

## Context

[ADR-001](001-runtime-talos-linux.md) selected **Talos Linux** on QEMU VMs, and
[ADR-002](002-cluster-topology.md) fixed the topology at one control-plane and
two workers. During Phase 1 implementation on the actual host (macOS, Apple
Silicon) **both Talos provisioning paths failed**, reproducibly:

1. **Talos on QEMU** (`talosctl cluster create qemu`, macOS `vmnet-shared`).
   The VM booted (Talos userspace came up, NIC `enp0s5` present) and even obtained
   an IPv6 ULA — so L2 was healthy — but **DHCPv4 never received an offer**
   (`network.OperatorSpecController: DHCP request/renew failed … no matching
   response packet received`). The node was therefore unreachable on IPv4
   (`talosctl` waited on `10.5.0.2:50000`, then `192.168.64.2:50000`; both timed
   out, ARP `(incomplete)`). Changing the cluster CIDR to match the macOS vmnet
   default (`192.168.64.0/24`) did not help, so this is not a subnet issue.
2. **Talos on Docker** (`talosctl cluster create docker` on Colima).
   Containers started, but Talos `machined` failed fatally setting up its `/etc`
   overlay: `failed to set up /etc overlay: FSCONFIG_SET_FD failed: bad file
   descriptor`. Talos-in-Docker depends on the Docker host kernel; the Colima
   (Lima) kernel does not satisfy the overlay requirements, so the API never
   became usable.

Docker Desktop (which ships a different LinuxKit kernel) was not installed, and
installing it adds weight and licensing terms that are not justified here.

## Decision

Run the cluster as **k3s via k3d**, on the Docker runtime provided by **Colima**
(macOS Virtualization.Framework):

```bash
k3d cluster create homelab \
  --servers 1 --agents 2 \
  --port "80:80@loadbalancer" --port "443:443@loadbalancer" \
  --k3s-arg "--disable=traefik@server:0"
```

- **Topology preserved:** one control-plane server + two agents.
- **Bundled Traefik disabled** at bootstrap so ingress is installed declaratively
  through Argo CD, per [ADR-003](003-gitops-argocd.md) and
  [ADR-006](006-ingress-traefik.md).
- Ports 80/443 are mapped to the k3d load balancer for later public ingress.

This **supersedes ADR-001** (runtime) and adjusts the *provisioning* of ADR-002
(the node count and roles are unchanged; "node" now means a k3d node).

## Alternatives considered

### Keep fighting Talos on this macOS host
- Manually booting Talos VMs in UTM with **bridged** networking (real LAN DHCP)
  is the most likely way to get Talos working, but it replaces `talosctl`-managed
  provisioning with hand-built VMs and remains unproven on this host.
- **Deferred, not rejected:** it is the migration target if/when a Linux or
  bare-metal host is available (see Consequences).

### Docker Desktop + Talos-in-Docker
- A different kernel might satisfy Talos' overlay requirements.
- Rejected: heavyweight, separate licensing, and a large dependency for the sake
  of one runtime choice.

### k3s single-node in Colima
- Simplest and reliable, but a single node demonstrates no scheduling, roles, or
  multi-node failure modes — the main reason for not choosing the simplest path
  in ADR-001.

### k0s / MicroK8s
- Viable distributions, but no advantage over k3s on this host and less homelab
  mindshare.

## Consequences

- **Gained:** a working, multi-node Kubernetes cluster today; the entire platform
  story (Argo CD, LGTM, SLOs, security, backups, incident drill) is unaffected by
  the distribution change.
- **Accepted trade-off:** k3d nodes are **containers, not virtual machines** —
  the exact property ADR-001 used to reject k3d. The reason for accepting it now
  is that the alternative on this host is *no cluster at all*. Kubernetes
  behaviour (scheduling, taints, CRDs, RBAC) is unaffected; kernel/OS-level
  behaviour is the shared Colima kernel, so Talos-specific capabilities
  (immutable OS, `talosctl` upgrades) are unavailable.
- **Host dependency:** the cluster runs inside Colima; the Docker runtime
  (`unix://$HOME/.colima/default/docker.sock`) and Colima's VM must be running.
  `k3d cluster stop/start` manages the cluster lifecycle.
- **Talos remains the intended production target.** When a Linux/bare-metal host
  is available, the plan is to re-adopt Talos (QEMU or metal); this ADR is the
  documented migration path, and the Talos project skill is retained.
- **Changelog note:** the earlier Colima-vs-Talos choice ([ADR-001](001-runtime-talos-linux.md))
  is superseded by this ADR for the *current host*; the reasoning there still
  explains why Talos is preferred where it can run.
