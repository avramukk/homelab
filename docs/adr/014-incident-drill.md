# ADR-014: Deliberate incident drill and postmortem

## Status
Accepted

## Date
2026-09-16

## Context

A system that has never failed has never been tested. Claiming "production-grade"
requires evidence of *operating* under failure: detection, diagnosis, mitigation,
and learning. A showcase that only ever shows a green cluster demonstrates
nothing about reliability.

## Decision

Run a **deliberate, documented incident drill** as a first-class deliverable:

1. Inject a controlled failure (e.g. kill the demo service's database
   connectivity, exhaust a node, break an ingress route).
2. Confirm **detection by the alerting pipeline** (the alert actually fires and
   routes to Telegram).
3. Diagnose using **telemetry only** — logs, metrics, traces — not by reading
   source.
4. Mitigate using the **linked runbook**.
5. Write a **blameless postmortem** in [`docs/incident-reports/`](../incident-reports/)
   and update the runbook + alerts based on what was learned.

## Alternatives considered

### No deliberate failure
- Pros: nothing breaks on purpose.
- Cons: misses the point — no evidence of detection response or recovery.
- Rejected.

### Wait for real incidents only
- Pros: authentic.
- Cons: unpredictable, possibly never happens; the showcase would lack an incident
  narrative indefinitely.
- Rejected as the *only* source: real incidents are still documented when they
  occur.

### Chaos engineering platform (Litmus/Chaos Mesh) from day one
- Pros: automated, repeatable experiments.
- Cons: another heavy component to run and secure; premature for a single-host
  lab.
- Rejected for now: start with a scripted drill; adopt a platform if/when the
  cadence justifies it.

## Consequences

- **Gained:** demonstrable MTTA/MTTR, a real postmortem artifact, and feedback
  into alerts and runbooks — the loop that actually improves reliability.
- **Accepted:** the drill causes real downtime of the demo service; scheduled, not
  a surprise.
- **Culture:** postmortems are **blameless** — the fault is in the system
  (missing telemetry, unclear runbook), never the operator.
- **Recurring:** drills are repeatable; recurring drills become part of the
  operational calendar.
