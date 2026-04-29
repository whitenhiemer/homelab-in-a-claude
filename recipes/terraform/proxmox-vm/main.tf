terraform {
  required_version = ">= 1.5.0"

  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = "~> 0.66.0"
    }
  }
}

provider "proxmox" {
  endpoint  = var.proxmox_endpoint
  api_token = var.proxmox_api_token
  insecure  = var.proxmox_insecure

  ssh {
    agent    = true
    username = "root"
  }
}

resource "proxmox_virtual_environment_vm" "vm" {
  name      = var.name
  node_name = var.proxmox_node
  vm_id     = var.vm_id
  tags      = var.tags
  on_boot   = true

  description = var.description

  bios = "seabios"

  # Mount the ISO for initial OS installation
  cdrom {
    file_id = var.iso
  }

  cpu {
    cores = var.cpu_cores
    type  = "x86-64-v2-AES"
  }

  memory {
    dedicated = var.memory_mb
    floating  = var.memory_balloon_mb
  }

  disk {
    datastore_id = var.storage
    interface    = "scsi0"
    size         = var.disk_size_gb
    file_format  = "raw"
    ssd          = true
    discard      = "on"
  }

  network_device {
    bridge = var.network_bridge
    model  = "virtio"
  }

  agent {
    enabled = true
    timeout = "15s"
  }

  serial_device {}

  operating_system {
    type = "l26"
  }

  lifecycle {
    ignore_changes = [
      cdrom,
    ]
  }
}
