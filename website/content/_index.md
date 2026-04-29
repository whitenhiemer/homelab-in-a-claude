---
title: "Homelab in a Claude"
type: docs
---

# Homelab in a Claude

Tell Claude what you want. Claude builds your homelab.

From bare metal to a full self-hosted stack — Proxmox, Kubernetes, reverse proxy, SSO, media server, VPN, monitoring — guided step by step by an AI that knows what it's doing.

---

## How it works

**1. Install the MCP servers**

Clone the repo and build the tools that give Claude access to your infrastructure:

```bash
git clone https://github.com/whitenhiemer/homelab-in-a-claude
cd homelab-in-a-claude
make mcp-install
```

**2. Open Claude Code in the project directory**

```bash
cd homelab-in-a-claude
claude
```

**3. Describe your hardware**

Claude will ask what machines you have, what network setup you're running, and what you want to self-host. It proposes an architecture, waits for confirmation, then provisions everything.

**4. Claude does the heavy lifting**

Claude runs Terraform to create VMs and containers, Ansible to configure services, and kubectl to manage the Kubernetes cluster — pausing whenever it needs you to do something physical.

---

## What Claude can build

| Service | What it does |
|---|---|
| Proxmox VE | Hypervisor — runs all your VMs and containers |
| Talos Linux | Immutable Kubernetes OS |
| Traefik | Reverse proxy with automatic TLS via Let's Encrypt |
| Authentik | SSO — single sign-on for all your services |
| Plex / Jellyfin | Media server |
| Sonarr / Radarr / SABnzbd | Automated media management |
| Prometheus + Grafana | Monitoring and alerting |
| WireGuard | VPN for remote access |
| Home Assistant | Home automation |
| TrueNAS | NAS with ZFS storage |

A single-machine setup is fully supported — you don't need a cluster.

---

[Get started →]({{< relref "/docs/getting-started" >}})
