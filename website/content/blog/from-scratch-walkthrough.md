---
title: "From Bare Metal to Self-Hosted Stack: A Walkthrough"
date: 2026-03-18
---

This is a walkthrough of standing up a homelab from scratch using this project — one machine, Proxmox, and Claude doing the heavy lifting. By the end you'll have Traefik with automatic TLS, Authentik for SSO, and the ARR media stack running behind a VPN.

**Prerequisites:** A machine you can install Proxmox on, a domain in Cloudflare, and a Cloudflare API token with DNS edit permissions.

---

## 1. Install Proxmox

Download the Proxmox VE ISO and flash it to a USB drive. Boot the machine from it, run through the installer (set a static IP, note it down), and verify you can reach the web UI at `https://<your-ip>:8006`.

You don't need to do anything else in the Proxmox UI yet.

---

## 2. Create a Proxmox API token

In the Proxmox web UI:

1. **Datacenter → Permissions → API Tokens → Add**
2. User: `root@pam`, Token ID: `claude`, uncheck "Privilege Separation"
3. Copy the token secret — you only see it once

You'll pass this as `PROXMOX_API_TOKEN` when setting up the MCP servers.

---

## 3. Clone the repo and build the MCP servers

```bash
git clone https://github.com/whitenhiemer/homelab-in-a-claude
cd homelab-in-a-claude
make mcp-build
```

This compiles all six MCP server binaries into `mcp/bin/`. Then register them with Claude Code:

```bash
PROXMOX_HOST=192.168.1.10 \
PROXMOX_API_TOKEN="root@pam!claude=<secret>" \
CF_API_TOKEN="<cloudflare-token>" \
bash scripts/install-mcp.sh
```

Verify all six are registered:

```bash
claude mcp list
```

---

## 4. Open Claude Code and describe your setup

```bash
claude
```

Start a conversation:

> I have a single Proxmox node at 192.168.1.10. My domain is example.com and it's on Cloudflare. I want to set up: a Traefik reverse proxy LXC, an Authentik SSO LXC, and the ARR media stack (Sonarr, Radarr, Prowlarr, SABnzbd). I have a 2TB drive I want to use as NAS storage via TrueNAS in a VM.

Claude will ask a few clarifying questions — LAN subnet, whether you want Plex or Jellyfin, which services need SSO protection — and then propose a network layout with specific IPs for each container and VM.

Review the proposal. If the IPs conflict with something on your network, say so. Once you confirm, Claude generates the Terraform config.

---

## 5. Apply Terraform

Claude generates a `terraform.tfvars` based on your answers and runs:

```
terraform init
terraform plan
```

It shows you the plan — the containers and VMs it will create, their IPs, their resources. You approve it, then:

```
terraform apply
```

After a minute or two you have:
- A Traefik LXC at the IP Claude picked
- An Authentik LXC
- A Debian LXC for the ARR stack
- A TrueNAS VM

All powered off — Terraform created the containers but didn't configure them yet.

---

## 6. Deploy services with Ansible

Claude walks through the services in dependency order.

**Traefik first** — everything else needs the reverse proxy running before it's accessible:

> Deploy Traefik on the LXC you just created. My Cloudflare API token is in the CF_API_TOKEN env var. My ACME email is me@example.com.

Claude runs the traefik ansible recipe. When it's done, `https://traefik.example.com` loads the dashboard.

**Authentik next:**

> Deploy Authentik on the auth LXC. Same domain.

The playbook generates a PostgreSQL password and secret key, writes them to `/opt/authentik/.env`, and starts the Docker stack. Claude prints the URL for the initial setup wizard: `https://auth.example.com/if/flow/initial-setup/`.

You create your admin account there. Claude then prompts you through the Authentik UI steps to configure forward auth — creating the provider, the application, and the outpost. This part is manual because Authentik's setup wizard doesn't have a stable API path.

**ARR stack:**

> Deploy the ARR stack. Use gluetun with Mullvad as the VPN provider. My Mullvad account number is in MULLVAD_ACCOUNT_NUMBER.

Claude deploys the arr-stack recipe, wires gluetun for Mullvad, and confirms the containers are up. Sonarr, Radarr, Prowlarr, SABnzbd, and Overseerr are all reachable via their Traefik routes.

---

## 7. TrueNAS storage

> Configure TrueNAS on the VM. The admin password is in TRUENAS_PASSWORD. The data disk is sdb.

Claude calls the TrueNAS REST API to create the `tank` pool on `/dev/sdb`, creates the media and backups datasets, enables NFS, and exports the shares. It then walks you through adding the NFS storage to Proxmox so you can bind-mount media into the ARR LXC.

---

## 8. Test the full stack

At this point:

- `https://traefik.example.com` — Traefik dashboard (protected by Authentik)
- `https://auth.example.com` — Authentik admin UI
- `https://sonarr.example.com` — Sonarr (SSO login via Authentik)
- `https://radarr.example.com` — Radarr
- `https://overseerr.example.com` — Overseerr
- NFS media share mounted in the ARR LXC at `/media`

TLS is automatic — Let's Encrypt via Cloudflare DNS challenge, renewed by Traefik.

---

## What Claude actually did

Throughout this walkthrough, Claude:

- Used the **Proxmox MCP** to verify the node was reachable and check existing resources before planning
- Used the **Terraform MCP** to run init, plan, and apply
- Used the **SSH MCP** to verify containers came up and check service logs
- Used the **Ansible MCP** to run each playbook
- Used the **Cloudflare MCP** to confirm DNS records resolved before testing HTTPS

At no point did it touch anything without showing you the plan first. Every `terraform apply`, every ansible run, every SSH command was preceded by an explanation of what it was about to do.

The next post goes deeper on how the MCP servers work and how to extend them.
