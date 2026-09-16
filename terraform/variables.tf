variable "cloudflare_api_token" {
  description = "Cloudflare API token: Zone DNS Edit, Zone Read, Account Cloudflare Tunnel Edit."
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
    Subdomain → in-cluster Service served openly through the Cloudflare Tunnel.
    Creates a proxied CNAME plus a tunnel ingress rule.
    Example: { "status" = "http://uptime-kuma.status.svc.cluster.local:3001" }
  EOT
  type        = map(string)
  default     = {}
}

variable "private_hostnames" {
  description = <<-EOT
    Subdomain → tailnet-only access. Creates a **DNS-only A record** pointing at
    `tailscale_ip`, so the name resolves only for devices on the tailnet, where
    Traefik routes it by Host header. Deliberately NOT proxied and NOT in the tunnel.
    Example: { "grafana" = "kube-prometheus-stack-grafana.observability.svc.cluster.local" }
  EOT
  type        = map(string)
  default     = {}
}

variable "tailscale_ip" {
  description = "Tailscale (tailnet) IP of the cluster host; target of private hostnames."
  type        = string
  default     = ""
}

variable "reserved_hostnames" {
  description = "Subdomains that resolve through the tunnel but serve nothing (catch-all 404)."
  type        = set(string)
  default     = []
}
