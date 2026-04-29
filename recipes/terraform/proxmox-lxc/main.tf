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

resource "proxmox_virtual_environment_container" "lxc" {
  node_name   = var.proxmox_node
  vm_id       = var.vm_id
  description = var.description
  tags        = var.tags

  unprivileged  = true
  started       = true
  start_on_boot = true

  operating_system {
    template_file_id = var.template
    type             = "debian"
  }

  cpu {
    cores = var.cpu_cores
  }

  memory {
    dedicated = var.memory_mb
  }

  disk {
    datastore_id = var.storage
    size         = var.disk_size_gb
  }

  network_interface {
    name   = "eth0"
    bridge = var.network_bridge
  }

  initialization {
    hostname = var.hostname

    ip_config {
      ipv4 {
        address = "${var.ip_address}/${var.subnet_prefix}"
        gateway = var.gateway
      }
    }

    dns {
      servers = var.nameservers
    }

    user_account {
      keys = var.ssh_public_keys
    }
  }

  features {
    nesting = var.enable_nesting
  }

  lifecycle {
    ignore_changes = [
      initialization[0].dns[0].domain,
    ]
  }
}
