# Argo CD bootstrap

One-time, idempotent install of Argo CD into the current kube-context, followed
by applying the root Application (app-of-apps). See
[ADR-017](../../docs/adr/017-argocd-bootstrap.md).

## Run

```bash
export DOCKER_HOST="unix://$HOME/.colima/default/docker.sock"
kubectl config use-context k3d-homelab
./install-argocd.sh
```

Re-running is safe: `helm upgrade --install` converges to the committed values
and `root-app.yaml` is applied idempotently.

## Access (over Tailscale — no ingress yet)

```bash
kubectl -n argocd port-forward svc/argocd-server 8080:80
# https is terminated later by Traefik/Cloudflare; for now the server is plain HTTP
open http://localhost:8080
```

Initial admin password:

```bash
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath='{.data.password}' | base64 -d
```

Change it on first login and remove `argocd-initial-admin-secret` afterwards.

## Verify

```bash
kubectl -n argocd get deploy,pods
kubectl -n argocd get applications
kubectl -n argocd get application root -o jsonpath='{.status.sync.status} {.status.health.status}'
```

Expected once children exist: `Synced Healthy`.

## Files

| File | Purpose |
|---|---|
| `install-argocd.sh` | Idempotent `helm upgrade --install` + `kubectl apply root-app.yaml` |
| `argocd-values.yaml` | Resource limits and `server.insecure` (TLS is terminated upstream) |
| `root-app.yaml` | App-of-apps root pointing at `infra/apps` |
