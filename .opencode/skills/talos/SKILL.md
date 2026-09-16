---
name: talos
description: Operate a Talos Linux Kubernetes cluster created with `talosctl cluster create` (QEMU VMs) — bootstrap, kubeconfig, health checks, upgrades, etcd operations, and cluster teardown. Use when working with Talos nodes, `talosctl`, `talosconfig`, machine config, `talosctl upgrade`, `talosctl health`, etcd snapshots, or when a Talos node or the cluster is unhealthy.
---

# Talos Linux

Talos has **no shell and no SSH**. Every node operation goes through the Talos
API via `talosctl`. Do not expect `ssh` or `kubectl` to work for node-level
actions.

## Local lab (this repo)

VMs run on the host `lab-host` via QEMU. Control-plane + workers are ephemeral
`talosctl cluster create` machines.

```bash
# Create the cluster: 1 control-plane + 2 workers (~4 GiB each)
talosctl cluster create \
  --name homelab \
  --controlplanes 1 \
  --workers 2 \
  --cpus 2 --memory 4096 --disk 20480

# Always confirm flags against your version:
talosctl cluster create --help
```

Requires `qemu-system-aarch64` on the host (brew) and a downloaded Talos image.

## Daily operations

```bash
talosctl -n <node-ip> health                  # cluster/node health
talosctl -n <node-ip> version                 # running Talos version
talosctl -n <node-ip> dmesg                   # kernel log
talosctl -n <node-ip> logs kubelet            # service log
talosctl -n <node-ip> dashboard               # interactive monitor
talosctl kubeconfig ./kubeconfig              # write kubeconfig
talosctl -n <node-ip> get members             # etcd members
```

## Config

Machine config is generated at create time into a temp dir; for a persistent
lab, export and version-control it **without secrets**.

```bash
talosctl gen config homelab https://<endpoint>:6443   # generates controlplane.yaml, worker.yaml, talosconfig
talosctl -n <node-ip> apply-config --file controlplane.yaml
```

Never commit `talosconfig` or machine secrets ([ADR-005](../../docs/adr/005-secrets-sealed-secrets.md)).

## Upgrade (manual — see runbook)

```bash
# 1. Inspect first
talosctl -n <node-ip> version
talosctl -n <node-ip> upgrade --dry-run --image ghcr.io/siderolabs/installer:<version>
# 2. Upgrade one node at a time (control-plane last is wrong — upgrade CP first, then workers)
talosctl -n <node-ip> upgrade --image ghcr.io/siderolabs/installer:<version>
# 3. Verify
talosctl -n <node-ip> version
kubectl get nodes -o wide
```

Full procedure: [`docs/runbooks/talos-upgrade.md`](../../docs/runbooks/).

## etcd

```bash
talosctl -n <cp-ip> etcd members
talosctl -n <cp-ip> etcd status
talosctl -n <cp-ip> etcd snapshot /tmp/etcd.snapshot
```

## Teardown

```bash
talosctl cluster destroy --name homelab
```

## Guardrails

- Verify the target node before any `upgrade`, `reset`, or `apply-config`.
- `talosctl reset` wipes a node — treat it like `format`.
- Node changes are not reconciled by Git; document any manual step in a runbook.
