output "zone_id" {
  description = "Resolved Cloudflare zone id."
  value       = data.cloudflare_zone.this.id
}

output "tunnel_id" {
  description = "Cloudflare Tunnel id (used in the CNAME target)."
  value       = cloudflare_zero_trust_tunnel_cloudflared.homelab.id
}

output "public_hostnames" {
  description = "Hostnames published through the tunnel."
  value       = keys(var.public_hostnames)
}
