# Security Hardening

> Status: `Planned`. See [ADR-012](../adr/012-security-baseline.md).

## Threat model (short)

| Asset | Threat | Control |
|---|---|---|
| Public dashboards | Data leak, defacement | read-only anonymous Grafana, no secrets in dashboards |
| Cluster API | Unauthenticated access | API reachable only over Tailscale; no public ports |
| Secrets in Git | Leak via public repo | Sealed Secrets ciphertext only; sealing key off-cluster |
| Container images | Vulnerable / malicious images | Trivy scan in CI; Kyverno image policies |
| Pod compromise | Lateral movement, privilege escalation | PSS `restricted`, non-root, read-only FS, dropped caps, NetworkPolicy |
| Supply chain | Tampered manifests | signed commits, reviewed PRs, Renovate + review |

## Baseline controls

- **Pod Security Standards:** `restricted` enforced at namespace level.
- **No root:** `runAsNonRoot`, `readOnlyRootFilesystem`, `allowPrivilegeEscalation: false`, all capabilities dropped.
- **Network:** Cilium NetworkPolicies, default-deny, explicit allow.
- **Images:** scanned with Trivy in CI; digests pinned where practical.
- **Admission:** Kyverno policies for baseline requirements (labels, resources, image provenance).
- **Secrets:** Sealed Secrets; no plaintext secret material in Git.
- **Public surface:** Cloudflare Tunnel only; dashboards read-only.

## Explicitly accepted risks

- **Grafana is private** ([ADR-020](../adr/020-grafana-private.md)). It is reachable
  only over Tailscale; anonymous access is disabled. Only the status page is
  public, and it discloses nothing beyond per-check up/down.
- **Single-disk local backups.** Loss of the host loses the backup repository;
  offsite is planned. Documented in [ADR-010](../adr/010-backup-restic-local.md).

## Verification

- `gitleaks git .` — no secrets in history or tree
- `kubectl get ns -L pod-security.kubernetes.io/enforce`
- `kubectl get clusterpolicy` (Kyverno)
- Trivy report attached to the relevant CI run
- Public dashboards reviewed for sensitive data before each release
