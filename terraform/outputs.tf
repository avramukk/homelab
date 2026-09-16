output "zone_id" {
  description = "Resolved Cloudflare zone id."
  value       = data.cloudflare_zone.this.id
}

output "tunnel_id" {
  description = "Cloudflare Tunnel id (CNAME target is <tunnel_id>.cfargotunnel.com)."
  value       = cloudflare_zero_trust_tunnel_cloudflared.homelab.id
}

output "public_hostnames" {
  description = "Hostnames served openly through the tunnel."
  value       = var.public_hostnames
}

output "private_hostnames" {
  description = "Hostnames resolvable only on the tailnet (A → tailscale_ip)."
  value       = var.private_hostnames
}

output "reserved_hostnames" {
  description = "Hostnames that resolve through the tunnel but serve nothing (404)."
  value       = var.reserved_hostnames
}
