# Architecture

## How the pieces fit together

```
User ──── Claude Code ──── CLAUDE.md (builder context)
                │
                ├── mcp/proxmox    ──── Proxmox API
                ├── mcp/terraform  ──── Terraform CLI + recipes/terraform/
                ├── mcp/ansible    ──── Ansible CLI + recipes/ansible/
                ├── mcp/kubectl    ──── kubectl + Helm
                ├── mcp/cloudflare ──── Cloudflare API
                └── mcp/ssh        ──── SSH to any host
```

Claude is the orchestrator. The MCP servers are its hands. Recipes are its playbook library.

## MCP servers

Each server is a standalone Go binary exposing a small set of tools via the MCP protocol. Claude calls these tools to act on the homelab. They are intentionally thin — no business logic, just execution with clean error output.

| Server | Tools |
|--------|-------|
| proxmox | list_nodes, list_vms, create_lxc, create_vm, delete_vm, get_vm_status |
| terraform | init, plan, apply, destroy, output |
| ansible | run_playbook, list_playbooks, check_inventory |
| kubectl | apply, get, delete, helm_install, helm_upgrade |
| cloudflare | list_records, create_record, delete_record |
| ssh | exec |

## Recipes

Parameterized templates Claude adapts to the user's environment. Each recipe is self-contained with its own variables file.

```
recipes/
  terraform/
    proxmox-lxc/     # Generic LXC container
    proxmox-vm/      # Generic VM
    talos-cluster/   # Full Talos K8s cluster
  ansible/
    traefik/         # Traefik reverse proxy
    authelia/        # SSO
    monitoring/      # Prometheus + Grafana
    arr-stack/       # Sonarr, Radarr, Prowlarr, qBittorrent
    wireguard/       # WireGuard VPN
    home-assistant/  # HAOS
    truenas/         # TrueNAS Scale
  k8s/
    metallb/         # MetalLB load balancer
    base-namespaces/ # Standard namespace setup
```

## Build phases

A typical build follows this order:

1. **Proxmox install** — user does this physically, Claude provides exact steps
2. **Network** — static IPs, DNS, Cloudflare DDNS if using a domain
3. **Infrastructure** — Terraform provisions VMs and LXCs on Proxmox
4. **Services** — Ansible deploys services into the containers
5. **Ingress** — Traefik + TLS routing
6. **Auth** — Authelia SSO (optional)
7. **Kubernetes** — Talos cluster bootstrap (optional, for users who want it)

Claude adapts this to the user's actual goals. Not every user needs every phase.

## Security model

- Credentials are never stored in files — Claude prompts for them at runtime as environment variables
- SSH uses key-based auth — Claude helps generate and distribute keys
- Secrets passed to Ansible use `--extra-vars` at runtime, not in playbook files
- Terraform state is local by default; Claude can help set up remote state if the user wants it
