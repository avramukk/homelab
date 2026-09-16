provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

data "cloudflare_zone" "this" {
  filter = {
    name = var.zone_name
  }
}

# The tunnel and its remote-managed configuration are fully owned by OpenTofu
# (ADR-019). cloudflared in the cluster joins it with the tunnel token.
resource "cloudflare_zero_trust_tunnel_cloudflared" "homelab" {
  account_id = data.cloudflare_zone.this.account.id
  name       = var.tunnel_name
  config_src = "cloudflare"
}

resource "cloudflare_zero_trust_tunnel_cloudflared_config" "homelab" {
  account_id = data.cloudflare_zone.this.account.id
  tunnel_id  = cloudflare_zero_trust_tunnel_cloudflared.homelab.id

  # Hostname → service rules are added when a service exists. Until then the
  # catch-all returns 404 for any hostname reaching the tunnel.
  config = {
    ingress = [
      { service = "http_status:404" }
    ]
  }
}

resource "cloudflare_dns_record" "public" {
  for_each = var.public_hostnames

  zone_id = data.cloudflare_zone.this.id
  name    = each.value
  type    = "CNAME"
  content = "${cloudflare_zero_trust_tunnel_cloudflared.homelab.id}.cfargotunnel.com"
  proxied = true
  ttl     = 1
  comment = "Managed by OpenTofu (homelab)"
}
