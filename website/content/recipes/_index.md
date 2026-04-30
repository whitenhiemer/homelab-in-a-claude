---
title: "Recipes"
weight: 2
bookFlatSection: true
---

# Recipes

Pre-built, parameterized templates that Claude uses to provision your homelab. Each recipe is a standalone Terraform module, Ansible playbook, or Kubernetes manifest.

Claude picks the right recipe for the job, fills in your specific IPs and hostnames, and applies it — showing you the plan before changing anything.

---

## Terraform

| Recipe | What it creates |
|---|---|
| [proxmox-lxc](https://github.com/whitenhiemer/homelab-in-a-claude/tree/master/recipes/terraform/proxmox-lxc) | Single Debian 12 LXC container with static IP and SSH key injection |
| [proxmox-vm](https://github.com/whitenhiemer/homelab-in-a-claude/tree/master/recipes/terraform/proxmox-vm) | Generic KVM VM, boots from ISO for OS installation |
| [talos-cluster](https://github.com/whitenhiemer/homelab-in-a-claude/tree/master/recipes/terraform/talos-cluster) | Talos Linux Kubernetes cluster VMs (control plane + workers) |

## Ansible

| Recipe | What it deploys |
|---|---|
| [traefik](https://github.com/whitenhiemer/homelab-in-a-claude/tree/master/recipes/ansible/traefik) | Traefik v3 reverse proxy with Let's Encrypt DNS-01 via Cloudflare |
| [monitoring](https://github.com/whitenhiemer/homelab-in-a-claude/tree/master/recipes/ansible/monitoring) | Prometheus + Grafana + Alertmanager + optional Proxmox PVE exporter |
| [arr-stack](https://github.com/whitenhiemer/homelab-in-a-claude/tree/master/recipes/ansible/arr-stack) | Sonarr, Radarr, Prowlarr, Bazarr, Overseerr, SABnzbd behind gluetun VPN |
| [wireguard](https://github.com/whitenhiemer/homelab-in-a-claude/tree/master/recipes/ansible/wireguard) | WireGuard VPN server with per-client keypair generation and QR codes |
| [authelia](https://github.com/whitenhiemer/homelab-in-a-claude/tree/master/recipes/ansible/authelia) | Authentik SSO via Docker Compose (PostgreSQL + Redis + server + worker) |
| [truenas](https://github.com/whitenhiemer/homelab-in-a-claude/tree/master/recipes/ansible/truenas) | TrueNAS Scale via REST API — ZFS pool, datasets, NFS shares |
| [home-assistant](https://github.com/whitenhiemer/homelab-in-a-claude/tree/master/recipes/ansible/home-assistant) | Traefik route for Home Assistant OS + trusted_proxies instructions |

## Kubernetes

| Recipe | What it applies |
|---|---|
| [base-namespaces](https://github.com/whitenhiemer/homelab-in-a-claude/tree/master/recipes/k8s/base-namespaces) | Base namespace set (ingress-system, apps, monitoring) |
| [metallb](https://github.com/whitenhiemer/homelab-in-a-claude/tree/master/recipes/k8s/metallb) | MetalLB IP pool for bare-metal LoadBalancer services |

---

## Using a recipe directly

Every recipe is usable without Claude — just copy the example vars file, edit it, and run:

```bash
# Terraform recipe
cd recipes/terraform/proxmox-lxc
cp terraform.tfvars.example terraform.tfvars
# edit terraform.tfvars
terraform init && terraform apply

# Ansible recipe
cd recipes/ansible/traefik
ansible-playbook playbook.yml -i inventory.yml \
  --extra-vars "cf_api_token=... acme_email=..."
```
