# ADR-013: Documentation and git-history conventions

## Status
Accepted

## Date
2026-09-16

## Context

This is a portfolio artifact. A visiting SRE should be able to land in the
repository and reconstruct **not just what was built, but why** — the sequence of
decisions, the alternatives rejected, and the operational procedures. That is a
property of the *history and docs*, not of the code alone.

A large, opaque "initial commit" or a stream of `fix`/`update` messages would
discard exactly the signal this repository exists to show.

## Decision

1. **ADRs are the backbone.** Every significant, expensive-to-reverse decision is
   an ADR in `docs/adr/`, numbered sequentially, **append-only**.
2. **One ADR = one commit:** `docs(adr): ADR-NNN <title>`.
3. **Linear history.** Conventional Commits, trunk-based, rebase (no merge
   commits).
4. **Phase tags** mark milestones (`v0.1.0-planning`, `v0.2.0-cluster`, …).
5. **Signed commits** (SSH signing) so integrity is verifiable on GitHub.
6. **PRs carry the decision record:** every PR states
   Decision / Alternatives / Validation / Rollback.
7. **Diagrams are mermaid** in-repo, rendered by GitHub and diffable.
8. **Docs change with the system,** in the same PR.

## Alternatives considered

### Squash everything into a single "init" commit
- Pros: tidy-looking history.
- Cons: destroys the decision narrative — the single most valuable thing here.
- Rejected.

### No ADRs; capture decisions in PR descriptions only
- Pros: lighter.
- Cons: PRs are not browsable as a decision log; context is scattered; no
  superseding chain.
- Rejected: ADRs plus PRs together give both the log and the discussion.

### External wiki (Confluence/Notion)
- Pros: easy editing.
- Cons: outside the repository; not versioned with the code; not reviewable in
  Git; a dead end for a public showcase.
- Rejected.

### Auto-generated docs from commits
- Pros: no manual effort.
- Cons: generated changelogs flatten *why*; they describe change, not reasoning.
- Rejected as the primary record.

## Consequences

- **Gained:** a repo that reads as a narrative; decisions are traceable to dates,
  alternatives, and consequences; history is verifiable (signed, linear).
- **Discipline cost:** every significant change needs an ADR and a clean commit;
  this is deliberate and slow by design.
- **Append-only** means superseding an ADR adds a new one — old reasoning is never
  erased, only contextualised.
- Enforced by convention and CI (docs validation), not by a gate that blocks work
  prematurely.
