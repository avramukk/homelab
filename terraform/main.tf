provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

data "cloudflare_zone" "this" {
  filter = {
    name = var.zone_name
  }
}

# Tunnel used by cloudflared running in the cluster (ADR-008).
resource "cloudflare_zero_trust_tunnel_cloudflared" "homelab" {
  account_id = data.cloudflare_zone.this.account.id
  name       = var.tunnel_name
  config_src = "cloudflare"
}

# Tunnel ingress + DNS per public hostname. Empty for now (DNS/tunnel only):
# each entry creates a proxied CNAME and a matching tunnel ingress rule.
resource "cloudflare_zero_trust_tunnel_cloudflared_config" "homelab" {
  count      = length(var.public_hostnames) > 0 ? 1 : 0
  account_id = data.cloudflare_zone.this.account.id
  tunnel_id  = cloudflare_zero_trust_tunnel_cloudflared.homelab.id

  config = {
    ingress = concat(
      [
        for host, service in var.public_hostnames : {
          hostname = "${host}.${var.zone_name}"
          service  = service
        }
      ],
      [
        # Catch-all required by cloudflared; must be last.
        { service = "http_status:404" }
      ]
    )
  }
}

resource "cloudflare_dns_record" "public" {
  for_each = var.public_hostnames

  zone_id = data.cloudflare_zone.this.id
  name    = each.key
  type    = "CNAME"
  content = "${cloudflare_zero_trust_tunnel_cloudflared.homelab.id}.cfargotunnel.com"
  proxied = true
  ttl     = 1
  comment = "Managed by OpenTofu (homelab)"
}
