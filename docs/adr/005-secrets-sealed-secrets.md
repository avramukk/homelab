# ADR-005: Secrets in Git — Sealed Secrets

## Status
Accepted

## Date
2026-09-16

## Context

GitOps makes the repository the source of truth, but the repository is **public**.
Workloads need secrets (database credentials, the Telegram bot token, the
Cloudflare tunnel token, Grafana admin credentials). These must reach the cluster
without ever being readable in Git history.

## Decision

Use **Sealed Secrets** (Bitnami `sealed-secrets` controller). Secrets are
encrypted client-side into `SealedSecret` resources that are safe to commit; only
the in-cluster controller can decrypt them.

## Alternatives considered

### SOPS + Age
- Pros: works for any file type (including `talosconfig` and non-Kubernetes
  material); git-diff-friendly; no in-cluster controller.
- Cons: requires age keys to be available to the decrypting step (Argo CD plugin
  or CI), which adds moving parts; key handling is manual.
- Rejected as the *primary* mechanism: Sealed Secrets keeps decryption entirely
  inside the cluster, which matches pull-based GitOps with less machinery.

### External Secrets Operator + Vault
- Pros: the strongest model — central secret store, dynamic credentials, audit.
- Cons: Vault is heavy and operationally demanding; running it on a single
  three-node host is overkill and would compete for RAM.
- Rejected: disproportionate to the scale.

### Secrets kept out of Git entirely
- Pros: simplest rule.
- Cons: breaks GitOps — the cluster state would no longer be reproducible from
  the repository; onboarding a reader would require undocumented manual steps.
- Rejected: contradicts [ADR-003](003-gitops-argocd.md).

## Consequences

- **Gained:** ciphertext is safe in a public repo; no external secret service;
  decryption happens in-cluster.
- **Critical dependency:** the controller's **sealing private key** is the only
  thing that can decrypt everything. It must be backed up **out of band and
  encrypted**. Losing it means all `SealedSecret`s must be re-sealed.
  This is a first-class backup item (see
  [ADR-010](010-backup-restic-local.md)).
- **Rotation:** rotating the sealing key requires re-sealing all secrets; the
  controller supports multiple keys during transition.
- **Scope:** Sealed Secrets only handles Kubernetes `Secret`s. Non-Kubernetes
  material such as `talosconfig` and Talos machine secrets stays out of the
  repository entirely ([ADR-001](001-runtime-talos-linux.md)).
