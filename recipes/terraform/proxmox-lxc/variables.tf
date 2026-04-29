variable "proxmox_endpoint" {
  description = "Proxmox API endpoint URL (e.g., https://192.168.86.29:8006)"
  type        = string
}

variable "proxmox_api_token" {
  description = "Proxmox API token: USER@REALM!TOKENID=SECRET"
  type        = string
  sensitive   = true
}

variable "proxmox_insecure" {
  description = "Skip TLS verification (set true for self-signed certs)"
  type        = bool
  default     = true
}

variable "proxmox_node" {
  description = "Proxmox node to create the LXC on"
  type        = string
}

variable "vm_id" {
  description = "Proxmox VM/CT ID (must be unique in the cluster)"
  type        = number
}

variable "hostname" {
  description = "Hostname for the container"
  type        = string
}

variable "description" {
  description = "Proxmox description for the container"
  type        = string
  default     = ""
}

variable "tags" {
  description = "Proxmox tags for the container"
  type        = list(string)
  default     = []
}

variable "template" {
  description = "Debian LXC template file ID (e.g., local:vztmpl/debian-12-standard_12.7-1_amd64.tar.zst)"
  type        = string
  default     = "local:vztmpl/debian-12-standard_12.7-1_amd64.tar.zst"
}

variable "cpu_cores" {
  description = "Number of CPU cores"
  type        = number
  default     = 1
}

variable "memory_mb" {
  description = "RAM in MB"
  type        = number
  default     = 512
}

variable "disk_size_gb" {
  description = "Root disk size in GB"
  type        = number
  default     = 8
}

variable "storage" {
  description = "Proxmox storage pool for the container disk"
  type        = string
  default     = "local-lvm"
}

variable "network_bridge" {
  description = "Proxmox network bridge"
  type        = string
  default     = "vmbr0"
}

variable "ip_address" {
  description = "Static IPv4 address (without prefix, e.g., 192.168.86.50)"
  type        = string
}

variable "subnet_prefix" {
  description = "Subnet prefix length (e.g., 24 for /24)"
  type        = number
  default     = 24
}

variable "gateway" {
  description = "Default gateway IP"
  type        = string
}

variable "nameservers" {
  description = "DNS nameservers"
  type        = list(string)
  default     = ["8.8.8.8", "8.8.4.4"]
}

variable "ssh_public_keys" {
  description = "SSH public keys to inject for Ansible/root access"
  type        = list(string)
  default     = []
}

variable "enable_nesting" {
  description = "Enable nesting (required for Docker inside LXC)"
  type        = bool
  default     = false
}
