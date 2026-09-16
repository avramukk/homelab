provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

data "cloudflare_zone" "this" {
  filter = {
    name = var.zone_name
  }
}

# DNS only. The tunnel itself is created out-of-band with the cloudflared CLI
# (`cloudflared tunnel create`), because that path authenticates with an origin
# certificate instead of an account-scoped API token.
#
# Tunnel ingress (which service a hostname reaches) is configured on the
# cloudflared side, not here.
resource "cloudflare_dns_record" "public" {
  for_each = var.public_hostnames

  zone_id = data.cloudflare_zone.this.id
  name    = each.value
  type    = "CNAME"
  content = "${var.tunnel_id}.cfargotunnel.com"
  proxied = true
  ttl     = 1
  comment = "Managed by OpenTofu (homelab)"
}
