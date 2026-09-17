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

| Runbook | Used for | Status |
|---|---|---|
| [demo-slo.md](demo-slo.md) | demo error budget / latency / down | ✅ written |
| [cloudflare-tunnel.md](cloudflare-tunnel.md) | public hostname unreachable, DNS | ✅ written |
| [argocd-bootstrap.md](argocd-bootstrap.md) | Argo CD install / recovery | ✅ written |
| `restore.md` | restoring from Restic | ⏳ planned (with ADR-010 implementation) |
| `cluster-recovery.md` | cluster or node down / rebuild | ⏳ planned |
| `talos-upgrade.md` | upgrading the node OS | ⏳ planned (Talos is deferred) |

> Statuses here must match reality: a runbook listed as written must exist, and an
> alert must link a runbook that exists. Keep this table honest.
