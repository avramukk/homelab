---
name: restic-backup
description: Operate Restic backups for this homelab — initialise a repository, back up volumes and database dumps, verify with restic check, prune with retention, and perform a real restore. Use when configuring the backup CronJob, troubleshooting failed backups, running the monthly restore drill, adding an offsite (B2) target, or when BackupFailed fires.
---

# Restic Backups

Policy and RPO/RTO: [ADR-010](../../docs/adr/010-backup-restic-local.md) and
[`docs/operations/backup-restore.md`](../../docs/operations/backup-restore.md).

## Repository

Local target now; B2 is a planned second repository (same tool, new backend).

```bash
export RESTIC_REPOSITORY=/backups/restic
export RESTIC_PASSWORD_FILE=/run/secrets/restic-password   # never inline a password
restic init
```

## Backup

```bash
# dump the DB first, then snapshot both the dump and the volumes
pg_dump -Fc "$DATABASE_URL" -f /backup/db.dump
restic backup /backup /data --tag homelab

# verify integrity after every run
restic check
```

In-cluster shape: a CronJob mounts persistent volumes **read-only**, runs a
pre-backup `pg_dump` hook, then `restic backup`.

## Retention

```bash
restic forget --keep-daily 30 --keep-weekly 8 --keep-monthly 12 --prune
```

## Inspect

```bash
restic snapshots
restic snapshots --tag homelab --latest
restic ls latest
restic stats
```

## Restore (the part that must be tested)

```bash
# 1. Restore to a scratch location first — never straight over live data
restic restore latest --target /restore

# 2. Validate the dump
pg_restore --list /restore/backup/db.dump

# 3. Move into place during a maintenance window
```

Full procedure and drill record: [`docs/runbooks/restore.md`](../../docs/runbooks/).

## Offsite (planned)

```bash
export RESTIC_REPOSITORY=b2:bucket-name:/homelab
restic init
restic copy --from-repo /backups/restic   # or back up directly to B2
```

## Guardrails

- A backup that has never been restored is **not** a backup — run the monthly drill.
- `restic check` failure is a `page` alert ([ADR-009](../../docs/adr/009-alerting-telegram-slo.md)).
- The repository password and (later) B2 keys are secrets — sealed, never committed.
- Independent of this repository: the Sealed Secrets key and Talos secrets are
  backed up encrypted **out of band**.
