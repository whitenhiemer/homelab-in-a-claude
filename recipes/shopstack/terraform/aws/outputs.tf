output "public_ip" {
  description = "Elastic IP — set this as the A record in Cloudflare"
  value       = aws_eip.shopstack.public_ip
}

output "instance_id" {
  description = "EC2 instance ID"
  value       = aws_instance.shopstack.id
}

output "ansible_inventory_line" {
  description = "Paste this into your inventory.ini under [shopstack]"
  value       = "shopstack ansible_host=${aws_eip.shopstack.public_ip} ansible_user=admin ansible_ssh_private_key_file=~/.ssh/id_ansible"
}
