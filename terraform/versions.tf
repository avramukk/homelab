# OpenTofu project for external, non-Kubernetes resources (ADR-011).
#
# Scope: things Kubernetes/Argo CD cannot manage — Cloudflare DNS and the
# Cloudflare Tunnel. In-cluster resources stay with Argo CD; do not duplicate
# them here.
#
# Usage:
#   export PATH="/opt/homebrew/bin:$PATH"
#   cd terraform
#   tofu init
#   tofu plan        # token comes from terraform.tfvars (gitignored) or TF_VAR_cloudflare_api_token
#   tofu apply
#
# Secrets: `terraform.tfvars` and state files are gitignored. Never commit them.

terraform {
  required_version = ">= 1.9"

  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.0"
    }
  }
}
