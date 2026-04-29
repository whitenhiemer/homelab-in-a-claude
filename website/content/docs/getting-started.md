---
title: "Getting Started"
weight: 1
---

# Getting Started

## Prerequisites

### Local tools

You need these installed on the machine where you run Claude Code:

| Tool | Version | Install |
|---|---|---|
| Claude Code | latest | `npm install -g @anthropic-ai/claude-code` |
| Go | 1.22+ | [go.dev/dl](https://go.dev/dl/) |
| Terraform | 1.5+ | [developer.hashicorp.com](https://developer.hashicorp.com/terraform/install) |
| Ansible | 2.14+ | `pip install ansible` |
| kubectl | any | [kubernetes.io/docs](https://kubernetes.io/docs/tasks/tools/) |
| talosctl | match cluster | [talos.dev](https://www.talos.dev/latest/talos-guides/install/talosctl/) |

### Hardware

You need at least one physical or virtual machine to install Proxmox on. Minimum specs:

- **CPU:** 4 cores (Intel or AMD, x86-64)
- **RAM:** 8 GB (16 GB recommended)
- **Disk:** 100 GB SSD
- **NIC:** 1 Gbps ethernet

A single node is fine. Multi-node clusters add HA but aren't required.

### Network

- A router that lets you assign static IPs (or DHCP reservations)
- Ideally a `/24` subnet like `192.168.86.0/24`
- A domain name if you want HTTPS — Cloudflare's free tier works well

---

## Step 1: Install Proxmox

Proxmox VE is the hypervisor that runs all your VMs and containers. Claude can't do this part — it requires physical access.

1. Download the Proxmox VE ISO from [proxmox.com/downloads](https://www.proxmox.com/en/downloads)
2. Flash it to a USB drive (`dd` or Balena Etcher)
3. Boot from USB and follow the installer
   - Set a static IP during install (e.g. `192.168.86.29`)
   - Note the root password — you'll need it for the API token
4. After reboot, access the web UI at `https://<ip>:8006`

{{< hint info >}}
Proxmox uses a self-signed TLS cert by default. Accept the browser warning.
{{< /hint >}}

---

## Step 2: Create a Proxmox API token

Claude uses the Proxmox API to create and manage VMs. Create a token with the right permissions:

1. Log into the Proxmox web UI
2. Go to **Datacenter → Permissions → API Tokens → Add**
3. User: `root@pam`, Token ID: `homelab`, uncheck "Privilege Separation"
4. Copy the token secret — it's only shown once

Your token will look like: `root@pam!homelab=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`

---

## Step 3: Install the MCP servers

```bash
git clone https://github.com/whitenhiemer/homelab-in-a-claude
cd homelab-in-a-claude

# Build all 6 MCP server binaries
make mcp-build

# Register them with Claude Code
PROXMOX_HOST=192.168.86.29 \
PROXMOX_API_TOKEN="root@pam!homelab=xxxx..." \
CF_API_TOKEN="your-cloudflare-token" \
make mcp-install
```

Verify they registered:

```bash
claude mcp list
```

You should see: `proxmox`, `terraform`, `ansible`, `kubectl`, `cloudflare`, `ssh`.

---

## Step 4: Start a session

```bash
cd homelab-in-a-claude
claude
```

Tell Claude what you have:

> "I have one machine with 16GB RAM and 500GB SSD running Proxmox at 192.168.86.29. I want to run a media server with Plex, automated downloads with Radarr and Sonarr, and a VPN. My domain is example.com managed on Cloudflare."

Claude will ask clarifying questions, propose an IP allocation plan and container layout, and wait for your approval before provisioning anything.

---

## What happens next

Claude works through phases in order:

1. **Download Proxmox templates** — Debian 12 LXC template for containers
2. **Terraform** — creates LXC containers with static IPs and SSH keys
3. **Ansible** — installs Docker and deploys each service
4. **DNS** — creates Cloudflare DNS records for your domain
5. **TLS** — Traefik gets Let's Encrypt certs via DNS-01 challenge

Each phase shows you what it's about to do and waits for confirmation before applying.
