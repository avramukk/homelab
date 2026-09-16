---
name: cilium
description: Install and operate Cilium as the CNI for this homelab — network policies (default-deny), connectivity testing, Hubble flow observation, and endpoint/identity debugging. Use when configuring NetworkPolicy, troubleshooting pod connectivity or DNS, verifying east-west isolation, enabling Hubble, or when traffic that should be blocked is allowed (or vice versa).
---

# Cilium

Cilium provides networking, `NetworkPolicy` enforcement, and (optionally) Hubble
observability. This cluster uses it for a default-deny posture
([ADR-012](../../docs/adr/012-security-baseline.md)).

## Install (GitOps-managed)

Declared in `infra/apps`, installed via Helm. For a manual smoke test only:

```bash
helm repo add cilium https://helm.cilium.io/
helm install cilium cilium/cilium --namespace kube-system \
  --set ipam.mode=kubernetes \
  --set kubeProxyReplacement=true \
  --set hubble.enabled=true \
  --set hubble.relay.enabled=true \
  --set hubble.ui.enabled=true
```

## Verify

```bash
cilium status --wait
cilium connectivity test          # full matrix; slow but thorough
kubectl -n kube-system get pods -l k8s-app=cilium
```

## Network policies

Default-deny per namespace, then explicit allow:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
  namespace: demo
spec:
  podSelector: {}
  policyTypes: [Ingress, Egress]
```

```bash
kubectl -n demo get networkpolicies
kubectl -n demo describe networkpolicy default-deny
```

## Observe traffic (Hubble)

```bash
cilium hubble port-forward &          # or use the Hubble UI
hubble observe --namespace demo
hubble observe --namespace demo --verdict DROPPED
hubble observe --to-pod demo/app-xxx --follow
```

## Debug connectivity

```bash
cilium endpoint list
cilium service list
kubectl -n kube-system exec ds/cilium -- cilium monitor --type drop
```

| Symptom | First checks |
|---|---|
| DNS broken | CoreDNS reachable under policy? `hubble observe --to-port 53` |
| One service unreachable | Missing allow rule for that namespace/pod |
| Everything blocked after rollout | Default-deny applied without allow rules |

## Guardrails

- Always test with `hubble observe ... --verdict DROPPED` before removing a policy.
- Keep allow rules as narrow as the workload requires — no `0.0.0.0/0` egress.
