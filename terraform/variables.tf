variable "cloudflare_api_token" {
  description = "Cloudflare API token: Zone → DNS → Edit, Zone → Zone → Read, Account → Cloudflare Tunnel → Edit."
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
  description = "Subdomains under zone_name published through the tunnel (DNS records)."
  type        = set(string)
  default     = []
}
