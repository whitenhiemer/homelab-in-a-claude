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
  description = "Default Proxmox node (used when per-VM node lists are empty)"
  type        = string
}

variable "cluster_name" {
  description = "Kubernetes cluster name (used in VM names and tags)"
  type        = string
  default     = "talos"
}

variable "talos_iso" {
  description = "Talos ISO file ID on Proxmox (e.g., local:iso/talos-amd64.iso)"
  type        = string
}

variable "storage" {
  description = "Proxmox storage pool for VM disks"
  type        = string
  default     = "local-lvm"
}

variable "network_bridge" {
  description = "Proxmox network bridge"
  type        = string
  default     = "vmbr0"
}

# --- Control Plane ---

variable "controlplane_count" {
  description = "Number of control plane nodes (1 for single-node, 3 for HA)"
  type        = number
  default     = 1
}

variable "controlplane_vmid_start" {
  description = "Starting VM ID for control plane nodes"
  type        = number
  default     = 400
}

variable "controlplane_nodes" {
  description = "Proxmox node for each control plane VM (index-matched, falls back to proxmox_node)"
  type        = list(string)
  default     = []
}

variable "controlplane_cores" {
  description = "CPU cores per control plane node"
  type        = number
  default     = 2
}

variable "controlplane_memory" {
  description = "RAM ceiling in MB per control plane node"
  type        = number
  default     = 4096
}

variable "controlplane_balloon" {
  description = "RAM floor in MB per control plane node (0 disables ballooning)"
  type        = number
  default     = 2048
}

variable "controlplane_disk_size_gb" {
  description = "OS disk size in GB per control plane node"
  type        = number
  default     = 50
}

variable "controlplane_ips" {
  description = "Static IPs to assign to control plane nodes (informational — Talos configures IPs)"
  type        = list(string)
}

# --- Workers ---

variable "worker_count" {
  description = "Number of worker nodes"
  type        = number
  default     = 2
}

variable "worker_vmid_start" {
  description = "Starting VM ID for worker nodes"
  type        = number
  default     = 410
}

variable "worker_nodes" {
  description = "Proxmox node for each worker VM (index-matched, falls back to proxmox_node)"
  type        = list(string)
  default     = []
}

variable "worker_cores" {
  description = "CPU cores per worker node"
  type        = number
  default     = 4
}

variable "worker_memory" {
  description = "RAM ceiling in MB per worker node"
  type        = number
  default     = 8192
}

variable "worker_balloon" {
  description = "RAM floor in MB per worker node (0 disables ballooning)"
  type        = number
  default     = 4096
}

variable "worker_disk_size_gb" {
  description = "OS disk size in GB per worker node"
  type        = number
  default     = 100
}

variable "worker_ips" {
  description = "Static IPs to assign to worker nodes (informational — Talos configures IPs)"
  type        = list(string)
}

variable "cluster_vip" {
  description = "Virtual IP for the Kubernetes API server (Talos VIP, must be in same subnet)"
  type        = string
}
