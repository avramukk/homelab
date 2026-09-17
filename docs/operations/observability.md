# Observability, SLOs and Alerting

> Status: `Partial`. Metrics/logs/traces and dashboards run; **alert delivery is
> not wired yet** (Alertmanager has a null receiver) and backup/cluster alerts are
> missing. See [ADR-007](../adr/007-observability-lgtm-otel.md) and
> [ADR-009](../adr/009-alerting-telegram-slo.md).

## Stack

- **Metrics:** Prometheus (RED + USE), long-term via Mimir if needed.
- **Logs:** Loki, structured JSON, correlated by `requestId`.
- **Traces:** Tempo, OpenTelemetry SDK in the demo service.
- **Collection:** Grafana Alloy / OTel Collector.
- **Dashboards:** Grafana, provisioned as code (no click-ops).

## Service Level Objectives

The demo service (`apps/demo`) is the SLO-bearing workload. Targets:

| SLI | SLO | Window |
|---|---|---|
| Availability (successful requests / total) | 99.5% | 30 days |
| Latency (p95 of `http_request_duration_seconds`) | < 300 ms | 30 days |
| Error rate (5xx / total) | < 1% | 30 days |

Error budget = 1 − SLO. Burn-rate alerts use a **multi-window** policy
(fast burn: 14.4× over 1h/5m; slow burn: 6× over 6h/30m) so pages are rare and
actionable.

## Alerting policy

Two severities only — a third tier becomes noise.

| Severity | Meaning | Route | Response |
|---|---|---|---|
| `page` | user-facing, act now | Telegram (immediate) | investigate via runbook |
| `ticket` | degradation, act this week | Telegram (digest) | schedule |

Rules:

1. Alert on **symptoms**, not causes (error rate, latency, availability — not CPU).
2. Every alert **links a runbook**.
3. Every alert has a threshold and duration justified by the SLO or history.
4. Every new alert is **test-fired once** before it is trusted.

## Dashboards

Grafana is **private** (Tailscale only, [ADR-020](../adr/020-grafana-private.md)).

| Dashboard | Notes |
|---|---|
| Service health (RED) | per-service rate/errors/duration |
| Cluster / nodes (USE) | capacity and saturation |
| SLO & error budget | burn-rate status |
| Operations | drill-down, dependency health |
