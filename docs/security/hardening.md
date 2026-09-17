# Security Hardening

> Status: `Partial`. Pod Security `restricted` is enforced on `demo` and `landing`
> only; several namespaces are unlabelled, and Trivy/Kyverno are not wired yet
> ([ADR-012](../adr/012-security-baseline.md) is being brought in line with
> reality).

## Threat model (short)

| Asset | Threat | Control |
|---|---|---|
| Public surface | Data leak, defacement | only the status page + landing are public; no secrets in them |
| Cluster API | Unauthenticated access | API reachable only over Tailscale |
| Secrets in Git | Leak via public repo | Sealed Secrets ciphertext only; sealing key off-cluster |
| Container images | Vulnerable / malicious images | multi-arch builds with provenance + SBOM (scanning: planned) |
| Pod compromise | Lateral movement, privilege escalation | PSS `restricted` (2 namespaces today), non-root, RO FS, dropped caps, NetworkPolicy (partial) |
| Supply chain | Tampered manifests | signed commits, reviewed PRs |

## Baseline controls — implemented vs planned

| Control | State |
|---|---|
| Pod Security Standards | `restricted` **enforced on `demo` + `landing`**; other namespaces unlabelled — *rollout in progress* |
| Non-root, RO rootfs, dropped caps, `allowPrivilegeEscalation: false` | applied to the workloads we author (demo, postgres, landing, cloudflared) |
| NetworkPolicy default-deny + explicit allow | `demo`, `landing` only; k3s **kube-router** enforces it (not Cilium) — *rollout in progress* |
| Sealed Secrets; no plaintext in Git | ✅ enforced (gitleaks pre-commit + CI) |
| Public surface: Cloudflare Tunnel only, tailnet for admin UIs | ✅ (Traefik private routers restricted to the tailnet CIDR) |
| Image scanning (Trivy) | ❌ not wired yet (CI builds with provenance/SBOM only) |
| Admission policies (Kyverno) | ❌ not installed yet |
| Image digest pinning | ❌ demo currently on a mutable tag |

## Explicitly accepted risks

- **Grafana is private** ([ADR-020](../adr/020-grafana-private.md)). It is reachable
  only over Tailscale; anonymous access is disabled. Only the status page is
  public, and it discloses nothing beyond per-check up/down.
- **Tailnet-only demo API.** `demo.avramukk.com` resolves only inside the tailnet
  (DNS A → tailnet IP, Traefik + TLS) — it is *not* public. Its data is
  disposable. Rate limiting / Access are unnecessary while it stays tailnet-only.
- **Single-disk local backups.** Loss of the host loses the backup repository;
  offsite is planned. Documented in [ADR-010](../adr/010-backup-restic-local.md).

## Verification

- `gitleaks git .` — no secrets in history or tree
- `kubectl get ns -L pod-security.kubernetes.io/enforce`
- `kubectl get clusterpolicy` (Kyverno)
- Trivy report attached to the relevant CI run
- Public dashboards reviewed for sensitive data before each release
