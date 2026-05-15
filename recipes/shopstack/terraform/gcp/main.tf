terraform {
  required_version = ">= 1.5.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.gcp_project
  region  = var.gcp_region
  zone    = var.gcp_zone
}

# Static external IP — no DDNS needed
resource "google_compute_address" "shopstack" {
  name   = "${var.client_name}-shopstack-ip"
  region = var.gcp_region
}

# Firewall rules
resource "google_compute_firewall" "shopstack_web" {
  name    = "${var.client_name}-shopstack-web"
  network = "default"

  allow { protocol = "tcp"; ports = ["80", "443"] }
  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["shopstack"]
}

resource "google_compute_firewall" "shopstack_mail" {
  name    = "${var.client_name}-shopstack-mail"
  network = "default"

  allow { protocol = "tcp"; ports = ["25", "465", "587", "143", "993", "995"] }
  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["shopstack"]
}

resource "google_compute_firewall" "shopstack_wireguard" {
  name    = "${var.client_name}-shopstack-wireguard"
  network = "default"

  allow { protocol = "udp"; ports = ["51820"] }
  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["shopstack"]
}

resource "google_compute_firewall" "shopstack_ssh" {
  name    = "${var.client_name}-shopstack-ssh"
  network = "default"

  allow { protocol = "tcp"; ports = ["22"] }
  source_ranges = var.admin_cidr_blocks
  target_tags   = ["shopstack"]
}

# Compute instance
resource "google_compute_instance" "shopstack" {
  name         = "${var.client_name}-shopstack"
  machine_type = var.machine_type
  zone         = var.gcp_zone
  tags         = ["shopstack"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-12"
      size  = var.disk_size_gb
      type  = "pd-balanced"
    }
  }

  network_interface {
    network = "default"
    access_config {
      nat_ip = google_compute_address.shopstack.address
    }
  }

  metadata = {
    ssh-keys  = "admin:${var.ssh_public_key}"
    user-data = <<-EOF
      #!/bin/bash
      hostnamectl set-hostname ${var.client_name}-shopstack
      apt-get update -q
      apt-get install -yq python3
    EOF
  }

  labels = { client = var.client_name, product = "shopstack" }
}
