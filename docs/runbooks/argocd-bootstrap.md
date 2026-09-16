# Runbook: Argo CD bootstrap / recovery

**Means:** Argo CD is not installed, is unreachable, or the `root` app-of-apps is
missing/broken and nothing reconciles.
**First check:** `kubectl -n argocd get pods` and
`kubectl -n argocd get application root`.
**Escalate to:** yourself — this is a single-operator lab; log the incident in
[`docs/incident-reports/`](../incident-reports/).

Related: [ADR-017](../adr/017-argocd-bootstrap.md), skill `.opencode/skills/argocd`.

## Prerequisites

```bash
export DOCKER_HOST="unix://$HOME/.colima/default/docker.sock"
kubectl config use-context k3d-homelab
kubectl cluster-info
```

## Install / re-install (idempotent)

```bash
git clone https://github.com/avramukk/homelab.git ~/homelab-bootstrap
cd ~/homelab-bootstrap
./infra/bootstrap/install-argocd.sh
```

## Verify

```bash
kubectl -n argocd get pods
kubectl -n argocd get applications
kubectl -n argocd get application root \
  -o jsonpath='{.status.sync.status} {.status.health.status}'
# expect: Synced Healthy
kubectl -n argocd logs deploy/argocd-repo-server --tail=50 | grep -i error
```

## Access (Tailscale only, no ingress yet)

```bash
kubectl -n argocd port-forward svc/argocd-server 8080:80
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath='{.data.password}' | base64 -d
```

## Recovery: root app deleted or drifted

```bash
kubectl apply -f infra/bootstrap/root-app.yaml
```

## Recovery: repo access fails

The repository is public, so no credentials are needed. If syncs fail with a
clone error, confirm the repo is reachable and the URL in `root-app.yaml`
matches, then check `argocd-repo-server` logs.

## Recovery: self-management wedged

If the `argocd` child Application (which manages Argo CD itself) is broken,
re-run `./infra/bootstrap/install-argocd.sh` — it converges the Helm release and
re-applies the root app, after which Git takes over again.

## Guardrails

- After bootstrap, change Argo CD **through Git**, not `helm upgrade`.
- Do not hand-edit resources Argo CD manages; the next sync reverts it.
- `--prune` deletes resources removed from Git — review the diff first.
