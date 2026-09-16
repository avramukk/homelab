# Incident Reports

Blameless postmortems for real incidents **and** for the deliberate chaos drill
required by [ADR-014](../adr/014-incident-drill.md). The point is to document
what happened, what we learned, and what changed — not who to blame.

## Template

```markdown
# <YYYY-MM-DD> — <short title>

## Summary
One paragraph a busy reader can act on.

## Impact
Who/what was affected, for how long, and how it was detected.

## Timeline
| Time (local) | Event |
|---|---|
| HH:MM | first signal |

## Root cause
The mechanism, not the person. Distinguish trigger vs cause.

## What went well
Signals that worked, shortcuts that paid off.

## What went badly
Gaps: missing telemetry, stale runbook, unclear ownership.

## Action items
| Action | Owner | Due | Status |
|---|---|---|---|

## Detection & response
Which alert fired (or which one should have), MTTA/MTTR.
```

## Index

| Date | Title | Type |
|---|---|---|
| — | (first drill pending) | chaos |
