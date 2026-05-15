output "public_ip" {
  description = "Static external IP — set this as the A record in Cloudflare"
  value       = google_compute_address.shopstack.address
}

output "instance_name" {
  description = "Compute instance name"
  value       = google_compute_instance.shopstack.name
}

output "ansible_inventory_line" {
  description = "Paste this into your inventory.ini under [shopstack]"
  value       = "shopstack ansible_host=${google_compute_address.shopstack.address} ansible_user=admin ansible_ssh_private_key_file=~/.ssh/id_ansible"
}
