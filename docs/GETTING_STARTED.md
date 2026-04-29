# Getting Started

## Prerequisites

- [Claude Code](https://claude.ai/code) with an Anthropic API key
- A machine to build on (bare metal or existing server — single node is fine)
- Basic comfort with the terminal

Local tools are installed as needed during the build. You don't need Terraform or Ansible pre-installed.

## Setup

### 1. Clone the repo

```bash
git clone https://github.com/whitenhiemer/homelab-in-a-claude.git
cd homelab-in-a-claude
```

### 2. Install the MCP servers

Each MCP server is a small Go binary. Build them all:

```bash
make mcp-build
```

### 3. Register with Claude Code

Add the MCP servers to your Claude Code config:

```bash
make mcp-install
```

This writes entries to `~/.claude/mcp.json` for each server.

### 4. Start building

Open Claude Code in this directory:

```bash
claude .
```

Then just tell Claude what you want:

> "I have two old ThinkCentre mini PCs and I want to run a media server and VPN."

Claude will take it from there.

## What Claude needs from you

- Answers to questions about your hardware and network
- Physical access to machines for the Proxmox install
- Credentials at runtime (Proxmox API token, Cloudflare API key, etc.) — Claude will ask for these when needed and never store them

## Single node vs. cluster

Single node is fully supported and a great starting point. You can always add nodes later — Claude knows how to expand the cluster without rebuilding from scratch.
