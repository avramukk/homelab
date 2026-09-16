variable "cloudflare_api_token" {
  description = "Cloudflare API token with DNS:Edit and Cloudflare Tunnel:Edit permissions for the zone."
  type        = string
  sensitive   = true
}

variable "zone_name" {
  description = "Cloudflare-managed apex zone."
  type        = string
  default     = "avramukk.com"
}

variable "tunnel_name" {
  description = "Name of the Cloudflare Tunnel used by cloudflared in the cluster."
  type        = string
  default     = "homelab"
}

variable "public_hostnames" {
  description = <<-EOT
    Hostname → in-cluster Service mappings published through the tunnel.
    Empty until a service actually exists (currently: DNS/tunnel only).
    Example: { "status" = "http://uptime-kuma.status.svc.cluster.local:80" }
  EOT
  type        = map(string)
  default     = {}
}
