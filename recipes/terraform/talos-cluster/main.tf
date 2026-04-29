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

# --- Control Plane Nodes ---

resource "proxmox_virtual_environment_vm" "controlplane" {
  count = var.controlplane_count

  name      = "${var.cluster_name}-cp-${count.index}"
  node_name = length(var.controlplane_nodes) > count.index ? var.controlplane_nodes[count.index] : var.proxmox_node
  vm_id     = var.controlplane_vmid_start + count.index
  tags      = ["kubernetes", "controlplane", var.cluster_name]
  on_boot   = true

  cdrom {
    file_id = var.talos_iso
  }

  cpu {
    cores = var.controlplane_cores
    type  = "x86-64-v2-AES"
    units = 1200
  }

  memory {
    dedicated = var.controlplane_memory
    floating  = var.controlplane_balloon
  }

  disk {
    datastore_id = var.storage
    interface    = "scsi0"
    size         = var.controlplane_disk_size_gb
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
    ignore_changes = [cdrom]
  }
}

# --- Worker Nodes ---

resource "proxmox_virtual_environment_vm" "worker" {
  count = var.worker_count

  name      = "${var.cluster_name}-worker-${count.index}"
  node_name = length(var.worker_nodes) > count.index ? var.worker_nodes[count.index] : var.proxmox_node
  vm_id     = var.worker_vmid_start + count.index
  tags      = ["kubernetes", "worker", var.cluster_name]
  on_boot   = true

  cdrom {
    file_id = var.talos_iso
  }

  cpu {
    cores = var.worker_cores
    type  = "x86-64-v2-AES"
    units = 1024
  }

  memory {
    dedicated = var.worker_memory
    floating  = var.worker_balloon
  }

  disk {
    datastore_id = var.storage
    interface    = "scsi0"
    size         = var.worker_disk_size_gb
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
    ignore_changes = [cdrom]
  }
}
