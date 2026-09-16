# Documentation

This directory is the **story of the homelab**: what was decided, why, and how it
is operated. It is written so that another SRE can walk in and follow the
reasoning end to end, without reading the code first.

## Read in this order

1. [`../README.md`](../README.md) — what this project is.
2. [`architecture/overview.md`](architecture/overview.md) — the system on one page.
3. [`adr/`](adr/) — every significant decision, in chronological order.
   This is the backbone: the history of *why*.
4. [`architecture/cluster.md`](architecture/cluster.md) — the concrete cluster.
5. [`operations/`](operations/) — credentials, SLOs, alerting/on-call, backups.
6. [`runbooks/`](runbooks/) — what to do when it breaks.
7. [`security/`](security/) — the hardening baseline.
8. [`incident-reports/`](incident-reports/) — blameless postmortems, incl. the
   deliberate chaos experiment.

## Conventions

- **ADRs are append-only.** A changed decision gets a new ADR that supersedes the
  old one; the old one is never deleted.
- **Docs change with the system**, in the same PR as the change.
- **Diagrams are mermaid** (`architecture/diagrams/*.mmd`) so they render on
  GitHub and stay reviewable in diffs.
- History is **linear** and follows [Conventional Commits](https://www.conventionalcommits.org/);
  one ADR per commit.

## Status legend

| Marker | Meaning |
|---|---|
| `Accepted` | decided and in effect |
| `Planned` | agreed, not implemented yet |
| `Superseded` | replaced by a later ADR |
