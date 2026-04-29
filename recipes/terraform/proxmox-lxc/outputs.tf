output "vm_id" {
  description = "Container ID"
  value       = proxmox_virtual_environment_container.lxc.vm_id
}

output "hostname" {
  description = "Container hostname"
  value       = var.hostname
}

output "ip_address" {
  description = "Container IP address"
  value       = var.ip_address
}
