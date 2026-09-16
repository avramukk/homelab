#!/usr/bin/env bash
# Bootstrap Argo CD into the current kube-context (idempotent), then hand
# control to the root Application (app-of-apps). See ../../docs/adr/017-argocd-bootstrap.md.
#
# Prereqs: helm, kubectl, a working kube-context, and (for the lab host) DOCKER_HOST
# pointing at the Colima socket.
set -euo pipefail

RELEASE="${ARGOCD_RELEASE:-argocd}"
NAMESPACE="${ARGOCD_NAMESPACE:-argocd}"
CHART="${ARGOCD_CHART:-argo/argo-cd}"
CHART_VERSION="${ARGOCD_CHART_VERSION:-}"
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

command -v helm >/dev/null || { echo "helm not found" >&2; exit 1; }
command -v kubectl >/dev/null || { echo "kubectl not found" >&2; exit 1; }
kubectl cluster-info >/dev/null || { echo "kube-context not reachable" >&2; exit 1; }

echo ">> adding argo helm repo"
helm repo add argo https://argoproj.github.io/argo-helm >/dev/null
helm repo update argo >/dev/null

echo ">> installing Argo CD ($CHART ${CHART_VERSION:-latest}) into $NAMESPACE"
# shellcheck disable=SC2086
helm upgrade --install "$RELEASE" "$CHART" \
  --namespace "$NAMESPACE" --create-namespace \
  ${CHART_VERSION:+--version "$CHART_VERSION"} \
  --values "$HERE/argocd-values.yaml" \
  --wait --timeout 10m

echo ">> applying root Application (app-of-apps)"
kubectl apply -f "$HERE/root-app.yaml"

echo ">> waiting for root Application to be created"
kubectl -n "$NAMESPACE" wait --for=condition=Established --timeout=60s crd/applications.argoproj.io >/dev/null

echo
echo "Argo CD is installed. Next:"
echo "  kubectl -n $NAMESPACE port-forward svc/argocd-server 8080:80"
echo "  kubectl -n $NAMESPACE get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d"
