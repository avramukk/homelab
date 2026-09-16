---
name: k3d
description: Operate this homelab's k3s cluster via k3d on Colima — create/stop/delete clusters, list nodes, merge kubeconfig, import images, access node containers, and troubleshoot. Use when working with k3d, k3s-in-Docker, `k3d cluster`, node containers, image imports, or when the local cluster is down or misbehaving.
---

# k3d (k3s in Docker)

The cluster runs **k3s inside Docker containers** managed by k3d, on the Colima
Docker runtime. Rationale and the Talos pivot: [ADR-016](../../docs/adr/016-runtime-k3s-k3d.md).

## Prerequisite

Colima must be running, and every command needs the Docker endpoint:

```bash
export DOCKER_HOST="unix://$HOME/.colima/default/docker.sock"
colima status
```

## Create the cluster

```bash
k3d cluster create homelab \
  --servers 1 --agents 2 \
  --port "80:80@loadbalancer" --port "443:443@loadbalancer" \
  --k3s-arg "--disable=traefik@server:0"
```

The bundled Traefik is disabled on purpose; ingress is installed via Argo CD
([ADR-006](../../docs/adr/006-ingress-traefik.md)).

## Daily operations

```bash
k3d cluster list
k3d cluster stop homelab          # pause node containers
k3d cluster start homelab
k3d node list
kubectl config use-context k3d-homelab
kubectl get nodes -o wide
```

## Load a local image (no registry)

```bash
docker build -t demo:dev ./apps/demo
k3d image import demo:dev -c homelab
```

## Node access (debug only)

```bash
docker exec -it k3d-homelab-server-0 sh
docker logs k3d-homelab-agent-0 --tail=100
```

## Delete / recreate (destructive)

```bash
k3d cluster delete homelab        # removes cluster + volumes; not Git-tracked
```

Recreate with the create command above. Desired state that matters lives in Git,
not in the cluster.

## Troubleshooting

| Symptom | First checks |
|---|---|
| `Cannot connect to the Docker daemon` | `colima status`; is `DOCKER_HOST` set? |
| Node `NotReady` | `docker ps`; `docker logs k3d-homelab-agent-0` |
| Ports 80/443 unreachable from the host | k3d loadbalancer container up? Colima port forwarding? |
| `kubectl` targets the wrong cluster | `kubectl config use-context k3d-homelab` |

## Guardrails

- **Nodes are containers** sharing the Colima kernel — no OS-level isolation to debug.
- Manage add-ons through Argo CD, not by hand; the next sync will revert manual edits.
- The cluster is disposable: `k3d cluster delete` + recreate must be enough to
  reproduce the platform, because Git holds the desired state.
