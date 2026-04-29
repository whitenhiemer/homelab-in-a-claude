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
  description = "Proxmox node to create the VM on"
  type        = string
}

variable "vm_id" {
  description = "Proxmox VM ID (must be unique in the cluster)"
  type        = number
}

variable "name" {
  description = "VM name shown in Proxmox UI"
  type        = string
}

variable "description" {
  description = "Proxmox description for the VM"
  type        = string
  default     = ""
}

variable "tags" {
  description = "Proxmox tags for the VM"
  type        = list(string)
  default     = []
}

variable "iso" {
  description = "ISO file ID for initial OS installation (e.g., local:iso/debian-12.iso)"
  type        = string
}

variable "cpu_cores" {
  description = "Number of CPU cores"
  type        = number
  default     = 2
}

variable "memory_mb" {
  description = "Maximum RAM in MB (balloon ceiling)"
  type        = number
  default     = 2048
}

variable "memory_balloon_mb" {
  description = "Minimum guaranteed RAM in MB (balloon floor, 0 disables ballooning)"
  type        = number
  default     = 0
}

variable "disk_size_gb" {
  description = "OS disk size in GB"
  type        = number
  default     = 32
}

variable "storage" {
  description = "Proxmox storage pool for the VM disk"
  type        = string
  default     = "local-lvm"
}

variable "network_bridge" {
  description = "Proxmox network bridge"
  type        = string
  default     = "vmbr0"
}
