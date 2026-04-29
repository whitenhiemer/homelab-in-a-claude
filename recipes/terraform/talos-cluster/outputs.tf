output "controlplane_vm_ids" {
  description = "VM IDs of control plane nodes"
  value       = proxmox_virtual_environment_vm.controlplane[*].vm_id
}

output "worker_vm_ids" {
  description = "VM IDs of worker nodes"
  value       = proxmox_virtual_environment_vm.worker[*].vm_id
}

output "controlplane_ips" {
  description = "Control plane IPs (from variables — Talos is configured separately)"
  value       = var.controlplane_ips
}

output "worker_ips" {
  description = "Worker IPs (from variables — Talos is configured separately)"
  value       = var.worker_ips
}

output "cluster_endpoint" {
  description = "Kubernetes API endpoint (via VIP)"
  value       = "https://${var.cluster_vip}:6443"
}

output "bootstrap_hint" {
  description = "Next step after apply"
  value       = "Run: talosctl apply-config --insecure --nodes <cp-ip> --file controlplane.yaml"
}
