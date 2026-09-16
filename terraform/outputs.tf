output "zone_id" {
  description = "Resolved Cloudflare zone id."
  value       = data.cloudflare_zone.this.id
}

output "tunnel_id" {
  description = "Cloudflare Tunnel id (CNAME target is <tunnel_id>.cfargotunnel.com)."
  value       = cloudflare_zero_trust_tunnel_cloudflared.homelab.id
}

output "public_hostnames" {
  description = "Hostnames published through the tunnel."
  value       = var.public_hostnames
}
