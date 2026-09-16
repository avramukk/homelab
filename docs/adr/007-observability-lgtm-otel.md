# ADR-007: Observability — self-hosted LGTM + OpenTelemetry

## Status
Accepted

## Date
2026-09-16

## Context

Observability is the centrepiece of an SRE showcase: metrics, logs, and traces,
plus dashboards that a visitor can actually look at. The lab could either
self-host the stack or pay for a managed backend.

Constraints:

- Must demonstrate *operating* the telemetry stack, not just consuming one.
- Bounded by the same 24 GiB host as everything else.
- Instrumentation should be vendor-neutral.

## Decision

Self-host the **LGTM stack** inside the cluster:

- **Prometheus** — metrics (RED for services, USE for resources)
- **Loki** — logs (structured JSON, correlated by `requestId`)
- **Grafana** — dashboards, provisioned as code
- **Tempo** — traces
- **OpenTelemetry** (SDK + Collector / Grafana Alloy) — instrumentation and
  pipeline

## Alternatives considered

### Grafana Cloud (managed)
- Pros: fast to adopt; the operator already has professional familiarity; no
  infrastructure to run; generous free tier.
- Cons: demonstrates *using* observability, not *running* it; data leaves the lab;
  the cluster carries less of the workload.
- Rejected as the primary backend: weaker showcase, and it removes the
  self-hosting story that this lab exists to tell.

### VictoriaMetrics instead of Prometheus
- Pros: single binary, excellent compression, very homelab-friendly.
- Cons: narrower market recognition; the goal is to demonstrate the standard.
- Rejected: a valid optimisation, not the headline.

### Datadog / New Relic
- Pros: industry-standard SaaS.
- Cons: paid, vendor lock-in, and again — nothing self-hosted to show.
- Rejected.

## Consequences

- **Gained:** full control of the telemetry pipeline; a real operational surface
  (retention, cardinality, scraping) that can be discussed in depth; OTel keeps
  instrumentation portable.
- **Accepted cost:** the stack consumes meaningful RAM on a 3-node, ~12 GiB
  cluster. Mitigations: conservative scrape intervals, bounded retention,
  histogram discipline to control cardinality.
- **Cardinality is the failure mode.** Labels come from small fixed sets (route
  templates, status classes), never user IDs or raw URLs. Percentiles
  (p50/p95/p99) are tracked, not averages.
- **Alerting** built on these signals is specified in
  [ADR-009](009-alerting-telegram-slo.md).
