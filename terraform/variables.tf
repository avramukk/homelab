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
  description = <<-EOT
    Subdomain → in-cluster Service mapping actually served through the tunnel.
    Creates both the proxied CNAME and the tunnel ingress rule.
    Keep this minimal: only services that are meant to be public.
    Example: { "status" = "http://uptime-kuma.status.svc.cluster.local:3001" }
  EOT
  type        = map(string)
  default     = {}
}

variable "reserved_hostnames" {
  description = <<-EOT
    Subdomains that resolve through the tunnel but serve no service (the tunnel
    catch-all answers 404). Use to reserve a name before its service exists.
  EOT
  type        = set(string)
  default     = []
}
