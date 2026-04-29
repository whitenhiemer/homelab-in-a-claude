---
title: "Contributing"
weight: 10
---

# Contributing

## Ways to contribute

- **Add a recipe** — a new Terraform module, Ansible playbook, or Kubernetes manifest
- **Improve an MCP server** — add tools, fix bugs, improve error messages
- **Write docs** — improve this site, add examples, document edge cases
- **Report issues** — open a GitHub issue if something doesn't work

## Repo structure

```
homelab-in-a-claude/
├── mcp/                  # MCP server binaries (Go)
│   ├── proxmox/          # Proxmox API tools
│   ├── terraform/        # Terraform runner
│   ├── ansible/          # Ansible runner
│   ├── kubectl/          # kubectl + Helm tools
│   ├── cloudflare/       # Cloudflare DNS tools
│   └── ssh/              # SSH exec tool
├── recipes/              # Infrastructure templates
│   ├── terraform/        # Terraform modules
│   ├── ansible/          # Ansible playbooks
│   └── k8s/              # Kubernetes manifests
├── scripts/              # install-mcp.sh and helpers
├── website/              # This site (Hugo)
├── CLAUDE.md             # Builder context for Claude
└── Makefile
```

## Adding a recipe

Recipes live in `recipes/<type>/<name>/`. Each recipe needs:

- The actual code (`main.tf` / `playbook.yml` / `manifest.yml`)
- A `README.md` explaining what it does, variables, and usage
- An example vars file if applicable

Keep recipes self-contained — they should work standalone without importing other recipes.

## Adding an MCP tool

MCP servers are Go programs in `mcp/<name>/`. To add a tool:

1. Write the handler function in the appropriate server
2. Register it with `mcpServer.AddTool(...)` in `main()`
3. Return errors via `errResult(err)` (sets `isError: true` per MCP spec)
4. Build with `make mcp-build` and smoke-test the `initialize` handshake

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' \
  | ./mcp/<name>/bin/<name>
```

## Building

```bash
# Build all MCP servers
make mcp-build

# Build the website
cd website && hugo server --buildDrafts

# Install MCP servers into Claude Code
make mcp-install
```

Go 1.22+ is required. If your system Go is a different version, set `GOROOT`:

```bash
GOROOT=/usr/lib/go-1.22 make mcp-build
```

## Commit style

- `feat(mcp): add proxmox_resize_disk tool`
- `fix(ssh): fall back to SSH agent when key file auth fails`
- `feat(recipes): add wireguard ansible playbook`
- `docs: add architecture overview`

Include `Co-Authored-By: Claude <noreply@anthropic.com>` when Claude wrote the code.

## Opening a PR

PRs go to the `master` branch. Keep them focused — one feature or fix per PR. The description should say what changed and why, not just what files were touched.
