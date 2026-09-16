# ADR-009: Alerting — Telegram, two severities, SLO burn-rate

## Status
Accepted

## Date
2026-09-16

## Context

Telemetry without alerting is a dashboard nobody watches. The lab needs an
alerting policy that demonstrates SRE practice: symptom-based alerts, SLOs with an
error budget, and routing that a solo operator can actually sustain without
developing alert fatigue.

There is no on-call rotation, no team, and no ticketing system — the "pager" is
one person's phone.

## Decision

- **Transport:** Alertmanager → **Telegram** (bot).
- **Severities:** exactly **two** — `page` (user-facing, act now) and `ticket`
  (degradation, act this week).
- **Alert basis:** **SLO burn-rate**, multi-window (fast: 14.4× over 1h/5m; slow:
  6× over 6h/30m), plus a small set of availability/liveness alerts.
- **Every alert links a runbook.**

## Alternatives considered

### Three or more severities (info/warning/critical)
- Pros: finer granularity; matches some enterprise runbooks.
- Cons: the middle tier becomes noise that trains the operator to ignore alerts.
- Rejected: two tiers force a real decision — "does a human need to act now, or
  later?".

### Email as the transport
- Pros: no bot, no token.
- Cons: slow to notice; easily buried; poor fit for "act now".
- Rejected.

### Discord
- Pros: good for a public showcase channel; rich embeds.
- Cons: chat-style reading is noisier for paging; Telegram is a better personal
  notification channel.
- Rejected as primary: Discord is a reasonable future secondary sink.

### No alerting (dashboards only)
- Pros: nothing to maintain.
- Cons: removes the entire incident-response story from the showcase.
- Rejected.

### Static threshold alerts on causes (CPU, memory)
- Pros: trivial to write.
- Cons: fire when nothing is wrong and miss unmodelled failures.
- Rejected: alerts fire on **symptoms users feel** (error rate, latency,
  availability); resource levels are for dashboards.

## Consequences

- **Gained:** low-noise paging; alerts that map directly to user impact; a
  demonstrable error-budget story for the public SLO dashboard.
- **Accepted:** the bot token is a secret (`SealedSecret`,
  [ADR-005](005-secrets-sealed-secrets.md)); Telegram is an external dependency.
- **Discipline:** each new alert must be test-fired once and link a runbook;
  noisy or non-actionable alerts are deleted, not ignored.
- **No escalation path** beyond one person — documented as an accepted limitation
  in [`docs/operations/observability.md`](../operations/observability.md).
