# Runbooks

Operational procedures, written to be followed under pressure. Every alert links
to the runbook that resolves it.

## Template

```markdown
# <Service> — <Issue>

**Means:** what is likely happening.
**First check:** the single command that narrows it down.
**Resolve:** ordered steps, copy-pasteable.
**Escalate to:** who to contact if this does not resolve it.
**Related:** dashboards, alerts, ADRs.
```

## Index

| Runbook | Used for | Linked alert |
|---|---|---|
| [talos-upgrade.md](talos-upgrade.md) | upgrading Talos / Kubernetes | — (planned maintenance) |
| [cluster-recovery.md](cluster-recovery.md) | cluster or node down / rebuild | `ClusterUnavailable` |
| [restore.md](restore.md) | restoring from Restic | `BackupFailed` |
| [demo-slo.md](demo-slo.md) | demo service error-budget / latency / down | `DemoErrorBudgetFastBurn`, `DemoErrorBudgetSlowBurn`, `DemoLatencyHigh`, `DemoDown` |
| [cloudflare-tunnel.md](cloudflare-tunnel.md) | public dashboard unreachable | `PublicDashboardDown` |
| [argocd-bootstrap.md](argocd-bootstrap.md) | re-bootstrapping GitOps | `ArgoCDOutOfSync` |

> Status: `Planned`. Runbooks are written as each subsystem lands, and updated
> during the incident drill ([ADR-014](../adr/014-incident-drill.md)).
