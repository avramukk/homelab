# ADR-012: Security hardening baseline

## Status
Accepted

## Date
2026-09-16

## Context

The lab is public-facing (dashboards) and hosts real data (the demo database). It
also claims to be "production-grade". A showcase that ignores hardening would
undercut its own premise, so security is treated as a first-class deliverable, not
a footnote.

## Decision

Adopt a **hardened baseline**:

1. **Pod Security** — `restricted` PSS enforced per namespace.
2. **Containers** — non-root, read-only root filesystem, all capabilities
   dropped, `allowPrivilegeEscalation: false`.
3. **Network** — **Cilium** CNI with default-deny `NetworkPolicy` and explicit
   allow rules; no flat east-west traffic.
4. **Images** — **Trivy** scanning in CI; digests pinned where practical.
5. **Admission** — **Kyverno** policies for baseline requirements (labels,
   resources, image provenance).
6. **Secrets** — Sealed Secrets ([ADR-005](005-secrets-sealed-secrets.md)); no
   plaintext in Git.
7. **Public surface** — outbound-only tunnel ([ADR-008](008-public-dashboard-cloudflare-tunnel.md)),
   anonymous dashboards are read-only and reviewed for sensitive data.

## Alternatives considered

### PSS + non-root + Trivy only (moderate)
- Pros: most of the benefit for less work.
- Cons: leaves east-west traffic unconstrained and admits unpoliced workloads;
  weaker story for a security-minded reviewer.
- Rejected: the incremental cost is small relative to the assurance gained.

### Defaults only (minimal)
- Pros: fastest.
- Cons: contradicts "production-grade"; no enforceable guarantees.
- Rejected.

### Full CIS benchmark / audited compliance
- Pros: maximum rigour.
- Cons: disproportionate for a single-host lab; much of CIS targets OS-level
  controls that Talos already provides or abstracts.
- Rejected: `restricted` PSS plus Cilium and policy admission covers the
  meaningful, demonstrable controls.

### Gatekeeper/OPA instead of Kyverno
- Pros: Rego is powerful and widely deployed.
- Cons: steeper learning curve; Kyverno's CEL-based policies are Kubernetes-native
  and (as of its graduation) the more common default.
- Rejected.

## Consequences

- **Gained:** a defensible security posture that can be walked through; several
  independently verifiable controls (`kubectl get ns -L pod-security...`,
  `kubectl get clusterpolicy`).
- **Accepted friction:** stricter policies reject naive manifests — a feature, but
  it slows first deploys.
- **Accepted risk — anonymous public Grafana:** dashboards expose operational
  metadata by design. A pre-release review ensures no secrets, credentials, or
  sensitive identifiers are shown.
- **Accepted risk — local-only backups** ([ADR-010](010-backup-restic-local.md)).
- Policy exceptions, where needed, must be narrow and documented in the relevant
  manifest.
