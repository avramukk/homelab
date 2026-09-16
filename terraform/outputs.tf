output "zone_id" {
  description = "Resolved Cloudflare zone id."
  value       = data.cloudflare_zone.this.id
}

output "public_hostnames" {
  description = "Hostnames published through the tunnel (DNS records managed here)."
  value       = var.public_hostnames
}

output "tunnel_target" {
  description = "CNAME target used for published hostnames."
  value       = "${var.tunnel_id}.cfargotunnel.com"
}
