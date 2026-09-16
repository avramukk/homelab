provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

data "cloudflare_zone" "this" {
  filter = {
    name = var.zone_name
  }
}

# The tunnel and its remote-managed configuration are owned by OpenTofu (ADR-019).
resource "cloudflare_zero_trust_tunnel_cloudflared" "homelab" {
  account_id = data.cloudflare_zone.this.account.id
  name       = var.tunnel_name
  config_src = "cloudflare"
}

locals {
  account_id  = data.cloudflare_zone.this.account.id
  tunnel_host = "${cloudflare_zero_trust_tunnel_cloudflared.homelab.id}.cfargotunnel.com"
  public_dns  = toset(concat(keys(var.public_hostnames), tolist(var.reserved_hostnames)))
  private_dns = toset(keys(var.private_hostnames))
}

# Public traffic: tunnel ingress + proxied CNAME.
resource "cloudflare_zero_trust_tunnel_cloudflared_config" "homelab" {
  account_id = local.account_id
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
  for_each = local.public_dns

  zone_id = data.cloudflare_zone.this.id
  name    = each.value
  type    = "CNAME"
  content = local.tunnel_host
  proxied = true
  ttl     = 1
  comment = "Public via Cloudflare Tunnel (OpenTofu)"
}

# Private traffic: DNS-only A records pointing at the tailnet address. A 100.x
# address is routable only inside the tailnet, so these names are unreachable
# from the internet while still being real domain names on your devices.
resource "cloudflare_dns_record" "private" {
  for_each = local.private_dns

  zone_id = data.cloudflare_zone.this.id
  name    = each.value
  type    = "A"
  content = var.tailscale_ip
  proxied = false
  ttl     = 300
  comment = "Private (Tailscale) — Traefik routes by host (OpenTofu)"
}
