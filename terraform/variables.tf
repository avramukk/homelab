variable "cloudflare_api_token" {
  description = "Cloudflare API token with Zone → DNS → Edit (and Zone → Zone → Read) for the zone."
  type        = string
  sensitive   = true
}

variable "zone_name" {
  description = "Cloudflare-managed apex zone."
  type        = string
  default     = "avramukk.com"
}

variable "tunnel_id" {
  description = <<-EOT
    Cloudflare Tunnel id created with the cloudflared CLI
    (`cloudflared tunnel create homelab`). DNS records point at
    "<tunnel_id>.cfargotunnel.com".
  EOT
  type        = string
}

variable "public_hostnames" {
  description = "Subdomains under zone_name published through the tunnel (DNS records only)."
  type        = set(string)
  default     = []
}
