# ADR-025: Run 1 is a pilot; a clean rebuild is planned

## Status
Accepted

## Date
2026-09-17

## Context

This cluster is the **first, deliberately experimental** attempt at the homelab
described in [`architecture/overview.md`](../architecture/overview.md). It was
built to find out where the sharp edges are before committing to a shape:

- which runtime actually works on the host (Apple Silicon, macOS, a Linux VM),
- how the public surface should be split between the tunnel and the tailnet,
- which components can genuinely be managed as code,
- where NetworkPolicy, secret handling and backups bite.

It is **not** intended to be the final system. A second run will be built from
scratch, reusing the decisions that held up and avoiding the ones that did not.

The evidence for that intent is already in the record: of 80 commits, 13 are
defect fixes, 2 are direct reverts, and 4 decisions had to be amended or
superseded within two days.

## Decision

Treat run 1 as a **pilot**, and make the documentation the deliverable:

1. **Document maximally.** Every significant decision is an ADR; every defect is
   an issue; every reversal is a *new* ADR that supersedes the old one. Nothing is
   silently corrected.
2. **Accept rough edges** that run 2 will remove — but only while they are written
   down, with evidence (commit, issue, ADR).
3. **Keep a running register** in [`lessons-learned.md`](../lessons-learned.md).
   Each entry states what hurt, the root cause, what run 2 should do, and links to
   the ADR/issue/commit that produced it.
4. **Do not reshape history.** Run 1's mistakes, reversals and dead ends are part
   of the record and stay visible.

## Alternatives considered

### Build run 1 as if it were the final system
- Would require committing to the runtime, the exposure policy and the component
  set before the unknowns are resolved.
- Rejected: those unknowns were discovered within hours (Talos blocked twice,
  Grafana published and reverted the same day), and unwinding a "final" system is
  far more expensive than unwinding a pilot.

### Iterate in place instead of rebuilding
- Cheaper short-term, and keeps the running cluster.
- Rejected: incremental repair carries the accumulated structural mistakes
  forward. A clean run 2 with a known-good baseline and a written list of pitfalls
  is cheaper than patching run 1 indefinitely.

### Document only the final state
- Produces a tidy repo with no history of *why*.
- Rejected: the value of run 1 is precisely the record of what went wrong.

## Consequences

- Run 1 is explicitly **not production**, despite being operated with
  production-grade practices (GitOps, SLOs, sealed secrets, audits).
- All defects found during the build are tracked as GitHub issues and closed with
  their commit evidence, so run 2 starts from a known pitfall list.
- Readers of this repository should read
  [`lessons-learned.md`](../lessons-learned.md) **before** the ADRs.
- Open work that run 1 never finished (alert delivery, backups, restore drill,
  policy/admission hardening) is tracked on the project board and is expected to
  be part of run 2's scope, not run 1's.
