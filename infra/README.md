# Infrastructure

GitOps root for the homelab. Everything the cluster runs is described here and
reconciled by Argo CD ([ADR-003](../docs/adr/003-gitops-argocd.md)).

```
infra/
├── bootstrap/   # one-time, imperative install of Argo CD + the root app
└── apps/        # app-of-apps: each file is a child Argo CD Application
```

## Bootstrap (once)

```bash
export DOCKER_HOST="unix://$HOME/.colima/default/docker.sock"
kubectl config use-context k3d-homelab
./infra/bootstrap/install-argocd.sh
```

See [bootstrap/README.md](bootstrap/README.md) for details and access
instructions. Rationale: [ADR-017](../docs/adr/017-argocd-bootstrap.md).

## After bootstrap

Argo CD manages itself and every workload from `infra/apps`. Change the cluster
by changing Git — not with `kubectl` or `helm`.
