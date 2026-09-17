# Documentation

This directory is the **story of the homelab**: what was decided, why, and how it
is operated. It is written so that another SRE can walk in and follow the
reasoning end to end, without reading the code first.

## Read in this order

1. [`../README.md`](../README.md) — what this project is.
2. [`lessons-learned.md`](lessons-learned.md) — what run 1 got wrong and what
   run 2 should change. **Read this first for context** ([ADR-025](adr/025-run-1-pilot.md):
   run 1 is a pilot).
3. [`architecture/overview.md`](architecture/overview.md) — the system on one page.
4. [`adr/`](adr/) — every significant decision, in chronological order.
   This is the backbone: the history of *why*.
5. [`architecture/cluster.md`](architecture/cluster.md) — the concrete cluster.
6. [`operations/`](operations/) — credentials, SLOs, alerting/on-call, backups.
7. [`runbooks/`](runbooks/) — what to do when it breaks.
8. [`security/`](security/) — the hardening baseline.
9. [`incident-reports/`](incident-reports/) — blameless postmortems, incl. the
   deliberate chaos experiment.
10. [`skills.md`](skills.md) — which agent skills this repository uses, and when.

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
