---
name: argocd
description: Operate Argo CD as the GitOps controller for this homelab — bootstrap the app-of-apps, sync/diff applications, debug OutOfSync or Degraded states, manage projects and repositories, and recover a lost Argo CD install. Use when working with Argo CD, `argocd` CLI, `Application`/`AppProject` resources, app-of-apps, sync policies, drift, or when a deployment is not reconciling.
---

# Argo CD

Argo CD pulls desired state from Git and reconcile it into the cluster. CI never
deploys — see [ADR-003](../../docs/adr/003-gitops-argocd.md).

## Repo conventions

- Root app points at `infra/apps` (app-of-apps).
- One directory per application; a `kustomization.yaml` per environment if needed.
- Secrets are `SealedSecret` resources, never plaintext ([ADR-005](../../docs/adr/005-secrets-sealed-secrets.md)).

## Access

```bash
kubectl -n argocd get applications
kubectl -n argocd get application <name> -o yaml
# initial admin password (change it immediately)
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d
```

```bash
export ARGOCD_SERVER=argocd.<domain>   # reachable over Tailscale only
argocd login $ARGOCD_SERVER
```

## Everyday operations

```bash
argocd app list
argocd app get <name>
argocd app diff <name>
argocd app sync <name> --prune
argocd app wait <name> --health --timeout 300
argocd app history <name>
argocd app rollback <name> <id>
```

## Debugging

| Symptom | First checks |
|---|---|
| `OutOfSync` | `argocd app diff <name>`; compare desired ref vs `targetRevision` |
| `Degraded` | `kubectl -n <ns> describe`; `argocd app get <name>` events |
| Sync fails on secrets | Sealed Secrets controller healthy? was ciphertext sealed for this cluster? |
| Nothing reconciles | `kubectl -n argocd logs deploy/argocd-application-controller` |

```bash
argocd app sync <name> --dry-run
kubectl -n argocd logs deploy/argocd-repo-server --tail=100
```

## Bootstrap / recovery

See [`docs/runbooks/argocd-bootstrap.md`](../../docs/runbooks/). The root
`Application` is the only thing applied by hand:

```bash
kubectl apply -f infra/bootstrap/root-app.yaml
```

## Guardrails

- `--prune` deletes resources no longer in Git — review the diff first.
- Do not hand-edit resources Argo CD manages; the next sync reverts it. Change Git.
- Prefer `syncPolicy.automated` with `prune` and `selfHeal` only for lower-risk apps.
