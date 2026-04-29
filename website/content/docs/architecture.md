---
title: "Architecture"
weight: 2
---

# Architecture

## Overview

The reference architecture this project was built from runs on a 5-node Proxmox cluster. A single-node setup works fine — the architecture scales down cleanly.

```
Internet
    │
    ▼
Cloudflare (DNS + optional proxy)
    │
    ▼
Router (port 80/443 → Traefik LXC)
    │
    ▼
Traefik LXC (.20)          ← single ingress for all HTTPS
    ├── media.example.com  → Plex LXC (.23)
    ├── sonarr.example.com → ARR stack LXC (.22)
    ├── grafana.example.com→ Monitoring LXC (.25)
    └── *.example.com      → ...
```

## Network layout

All services live in a single `/24` subnet. IPs are assigned statically via Terraform.

| Range | Purpose |
|---|---|
| `.1` | Router / gateway |
| `.20–.39` | LXC service containers |
| `.40–.49` | VMs (TrueNAS, Home Assistant) |
| `.100–.113` | Kubernetes (VIP + nodes) |
| `.150–.199` | MetalLB pool (K8s LoadBalancer IPs) |
| `.200+` | Proxmox nodes |

## Components

### Proxmox VE

The hypervisor. Runs all LXC containers and VMs. Claude communicates with it via the REST API — creating containers, starting/stopping VMs, checking resource usage.

Multi-node clusters use shared Ceph storage so VMs can live-migrate between nodes.

### LXC containers

Lightweight Linux containers — faster to create and cheaper on RAM than full VMs. All services except TrueNAS and Home Assistant run as LXC containers.

All containers are **unprivileged Debian 12**, with nesting enabled for the ones that run Docker.

### Traefik

The single reverse proxy for all HTTPS traffic. Configured with:
- HTTP → HTTPS redirect on port 80
- Let's Encrypt TLS via Cloudflare DNS-01 challenge (works behind NAT, supports wildcards)
- Per-service dynamic route configs in `/etc/traefik/dynamic/*.yml`
- Prometheus metrics endpoint for monitoring

### Authentik

SSO identity provider. Protects services that don't have their own auth (Grafana, Traefik dashboard, etc.) via Traefik's `forwardAuth` middleware. Users authenticate once and get access to everything.

### Kubernetes (optional)

Talos Linux — immutable, API-driven Kubernetes OS. Used for workloads that benefit from K8s (autoscaling, rolling deploys, Helm charts). Not required for a basic homelab.

MetalLB provides real IPs for `LoadBalancer` services on bare metal.

### TrueNAS

NAS VM with ZFS storage pools on dedicated physical disks. Serves NFS shares to:
- ARR stack (media library, downloads)
- Plex / Jellyfin (read-only media access)
- Home Assistant (backups)

### Monitoring

Prometheus + Grafana + Alertmanager on a dedicated LXC. Scrapes:
- Node exporter (host metrics)
- Prometheus PVE exporter (Proxmox cluster metrics)
- Traefik metrics (HTTP request rates, TLS cert expiry)
- Container metrics via cAdvisor

Alerts route to Discord webhooks.

## MCP servers

Claude interacts with the infrastructure through 6 MCP servers:

| Server | Tools | Purpose |
|---|---|---|
| `proxmox` | 11 tools | Create/manage VMs and LXC containers |
| `terraform` | 6 tools | Run Terraform in recipe directories |
| `ansible` | 4 tools | Run playbooks, ping hosts, list inventory |
| `kubectl` | 5 tools + 4 Helm | Kubernetes and Helm operations |
| `cloudflare` | 5 tools | Manage DNS records |
| `ssh` | 1 tool | Run commands on remote hosts |

All servers communicate with Claude Code over stdio using the MCP protocol. Credentials are passed as environment variables at registration time via `make mcp-install`.
