terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# Elastic IP — static public IP for DNS, no DDNS needed
resource "aws_eip" "shopstack" {
  domain = "vpc"
  tags   = { Name = "${var.client_name}-shopstack" }
}

resource "aws_eip_association" "shopstack" {
  instance_id   = aws_instance.shopstack.id
  allocation_id = aws_eip.shopstack.id
}

# Security group — ports required for ShopStack
resource "aws_security_group" "shopstack" {
  name        = "${var.client_name}-shopstack"
  description = "ShopStack: web, mail, WireGuard"

  # SSH (WireGuard handles remote mgmt; keep SSH restricted to your IP)
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = var.admin_cidr_blocks
  }

  # HTTP + HTTPS (Traefik)
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # SMTP
  ingress {
    from_port   = 25
    to_port     = 25
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # SMTP submission (TLS)
  ingress {
    from_port   = 465
    to_port     = 465
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 587
    to_port     = 587
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # IMAP + IMAPS
  ingress {
    from_port   = 143
    to_port     = 143
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 993
    to_port     = 993
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # POP3S
  ingress {
    from_port   = 995
    to_port     = 995
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # WireGuard remote management
  ingress {
    from_port   = 51820
    to_port     = 51820
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${var.client_name}-shopstack" }
}

# EC2 instance
data "aws_ami" "debian12" {
  most_recent = true
  owners      = ["136693071363"] # Debian official

  filter {
    name   = "name"
    values = ["debian-12-amd64-*"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }
}

resource "aws_instance" "shopstack" {
  ami                    = data.aws_ami.debian12.id
  instance_type          = var.instance_type
  key_name               = var.key_pair_name
  vpc_security_group_ids = [aws_security_group.shopstack.id]
  subnet_id              = var.subnet_id

  root_block_device {
    volume_type           = "gp3"
    volume_size           = var.disk_size_gb
    delete_on_termination = true
  }

  user_data = <<-EOF
    #!/bin/bash
    hostnamectl set-hostname ${var.client_name}-shopstack
    apt-get update -q
    apt-get install -yq python3
  EOF

  tags = { Name = "${var.client_name}-shopstack", Client = var.client_name }
}
