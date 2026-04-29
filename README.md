# Homelab in a Claude

An AI-guided homelab builder powered by Claude. Tell Claude what you want, answer a few questions about your hardware, and Claude builds your homelab with you — step by step, from bare Proxmox to a full-stack self-hosted environment.

## What it does

Claude acts as your homelab engineer. It asks about your hardware and goals, proposes an architecture, provisions infrastructure via Proxmox/Terraform, deploys services via Ansible, and tells you exactly what to do when human hands are required (cable runs, BIOS settings, ISO installs).

**What Claude can build:**
- Proxmox VE cluster setup
- Talos Linux Kubernetes cluster
- Reverse proxy + TLS (Traefik + Cloudflare)
- SSO (Authelia)
- Media stack (Plex, Jellyfin, *arr)
- Monitoring (Prometheus + Grafana)
- VPN (WireGuard)
- Home Assistant
- NAS (TrueNAS)
- ...and more via community recipes

## How it works

1. Install the MCP servers (see [Getting Started](docs/GETTING_STARTED.md))
2. Open Claude Code and point it at this repo
3. Tell Claude what you want to build
4. Claude takes it from there — asking questions, running tools, guiding you through the physical steps

## Repository layout

```
mcp/            # MCP servers — the tools Claude uses to act on your homelab
recipes/        # Terraform modules, Ansible playbooks, K8s manifests
website/        # Project website (Hugo)
docs/           # Architecture, contributing guide, getting started
CLAUDE.md       # System context that makes Claude a homelab expert
```

## Requirements

- A machine (or machines) to build on — bare metal or a single server works
- Claude Code with an Anthropic API key
- Local tools: `terraform`, `ansible`, `kubectl`, `talosctl` (installed as needed)

## Contributing

See [CONTRIBUTING.md](docs/CONTRIBUTING.md). Recipes and MCP server improvements welcome.

## License

MIT
