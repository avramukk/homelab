# ADR-010: Backups — in-cluster Restic to local disk

## Status
Accepted

## Date
2026-09-16

## Context

Data that cannot be restored does not exist. The cluster has persistent state that
is *not* reconstructible from Git: the demo database, application volumes, and the
Sealed Secrets sealing key (without which every committed `SealedSecret` is
unreadable).

Constraint: there is **no NAS and no offsite target today**. Only the host's local
disk is available without additional spend.

## Decision

Run **Restic as an in-cluster CronJob**, writing to a repository on the host's
local disk. Define and publish **RPO 24h / RTO ≤ 4h**. Run a **monthly restore
drill**. Treat offsite (Backblaze B2) as a **known limitation** to be addressed
later, not as already solved.

## Alternatives considered

### Borg
- Pros: strong compression, mature; excellent for local/NAS targets.
- Cons: requires a Borg server on the far side; less ergonomic with object storage
  when offsite is added.
- Rejected: Restic's "dumb storage" model makes the eventual B2 step trivial.

### Backups on the host (launchd/cron) instead of in-cluster
- Pros: independent of the cluster; survives a cluster outage.
- Cons: splits the operational surface — backup logic lives outside GitOps;
  harder to review and version.
- Rejected: keep the definition in Git and reconciled by Argo CD; the host-level
  copy of the sealing key is a separate, out-of-band item.

### Kubernetes-native (Velero) or replicated storage (Longhorn/Ceph)
- Pros: volume snapshots; replication.
- Cons: replication is meaningless on a single host; Velero assumes object storage
  and adds complexity; Ceph is drastically over-scoped.
- Rejected.

### No backups (Git is the source of truth)
- Pros: zero cost.
- Cons: false — Git cannot restore the database, volumes, or sealing key.
- Rejected.

## Consequences

- **Gained:** a defined, tested backup path with explicit RPO/RTO and a
  repeatable, GitOps-managed definition.
- **Accepted limitation — no offsite copy.** A local-disk repository does not
  survive disk failure or host loss; this is the single largest risk in the lab
  and is documented as such. Migrating the repository to B2 is a small change
  (same Restic repository, new backend).
- **Restore is proven, not assumed:** the monthly drill is recorded in
  [`docs/incident-reports/`](../incident-reports/).
- **`restic check`** runs after every backup; a failed check is a `page`
  (see [ADR-009](009-alerting-telegram-slo.md)).
- **Highest-value out-of-band items:** the Sealed Secrets sealing key and Talos
  machine secrets must be backed up encrypted, separately from this Restic
  repository.
