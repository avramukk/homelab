# Backups and Restore

> Status: `Planned`. See [ADR-010](../adr/010-backup-restic-local.md).

## Model

Restic runs as an **in-cluster CronJob**, reading persistent volumes and the
database dump, and writing to a repository on local disk. Offsite (Backblaze B2)
is a **known limitation** for now and is planned as a follow-up.

| Property | Target |
|---|---|
| RPO (max acceptable data loss) | 24 hours |
| RTO (max acceptable restore time) | 4 hours |
| Restore drill | monthly, documented in `incident-reports/` |
| Verification | `restic check` after every run |

## What is backed up

| Source | Method |
|---|---|
| Kubernetes manifests | Git (the source of truth — no backup needed) |
| Demo service database | `pg_dump` pre-hook, then Restic |
| Persistent volumes (apps) | Restic, volume mounted read-only |
| Sealed-secrets sealing key | **out of band**, encrypted, off-cluster |
| Talos machine secrets / `talosconfig` | **out of band**, encrypted, off-cluster |

> The sealing key and Talos secrets cannot be recovered from Git. Losing them
> means losing the ability to decrypt secrets / talk to the cluster. They are the
> highest-priority items to back up.

## Retention

```
--keep-daily 30 --keep-weekly 8 --keep-monthly 12 --prune
```

## Known limitations

- **No offsite copy yet.** A local-disk repository does not survive disk or host
  loss. Tracked in [ADR-010](../adr/010-backup-restic-local.md).
- Single-node storage means no replication; restore depends on the Restic
  repository being intact.

## Restore procedure

See [`runbooks/restore.md`](../runbooks/restore.md).
