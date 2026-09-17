# Runbook: demo service SLO alerts

Related: [ADR-009](../adr/009-alerting-telegram-slo.md),
[`docs/operations/observability.md`](../operations/observability.md).

> **Rollback is a Git operation.** Every application syncs with `selfHeal: true`,
> so `kubectl rollout undo` / manual edits are reverted by Argo CD on the next
> sync. To roll back, revert the commit (or use `argocd app rollback` for an
> emergency), then let Argo CD reconcile.

SLOs: availability ≥ 99.5% (error budget 0.5%), p95 < 300 ms.

## DemoErrorBudgetFastBurn (page)

**Means:** the service is burning its error budget ~14× too fast — a real,
user-visible failure, not noise.
**First check:**

```bash
kubectl -n demo get pods
kubectl -n demo logs deploy/demo --tail=50 | grep '"status":5'
```

**Likely causes and fixes**

| Cause | Check | Fix |
|---|---|---|
| Bad deploy | `kubectl -n demo rollout history deploy/demo` (read-only) | **`git revert <commit>`** — Argo CD re-syncs. Emergency: `argocd app rollback demo <id>` |
| Database unreachable | `kubectl -n demo logs deploy/demo | grep db_` | check `postgres-0`, NetworkPolicy, secret |
| Dependency erroring | traces in Tempo for the failing route | follow the error span to its cause |

A controlled way to reproduce it: `GET /api/error` (always 500).

## DemoErrorBudgetSlowBurn (ticket)

**Means:** a steady leak of the error budget. Not paging, but schedule it.
**First check:** the error ratio by route in Grafana (route label), and whether a
single route dominates.
**Fix:** treat as a normal defect — find the route, reproduce locally, fix, ship.

## DemoLatencyHigh (ticket)

**Means:** p95 above 300 ms over 10m.
**First check:** `demo:latency_p95:rate10m` and the per-route latency panel;
Postgres slow queries (`pg_stat_statements` if enabled).
Reproduce: `GET /api/slow?ms=400`.

## DemoDown (page)

**Means:** Prometheus can scrape no demo replica.
**First check:**

```bash
kubectl -n demo get pods
kubectl -n demo get servicemonitor demo -o yaml
kubectl -n observability logs statefulset/prometheus-kube-prometheus-stack-prometheus --tail=50
```

Common causes: all replicas failing readiness (then the DB check above), or the
ServiceMonitor selector/labels drifting from the Service labels.

## Verification after any fix

- Alert resolves and the error budget stops draining.
- `demo:error_ratio:rate5m` back below the threshold.
- Update this runbook if a step was wrong or missing.
