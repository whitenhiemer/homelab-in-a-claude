#!/usr/bin/env bash
# install-mcp.sh — register homelab-in-a-claude MCP servers with Claude Code.
#
# Run after `make mcp-build`. Each server binary must exist in mcp/<name>/bin/<name>.
# Servers requiring credentials read them from environment variables at runtime —
# set them here or export them in your shell profile.
#
# Usage:
#   CF_API_TOKEN=... PROXMOX_HOST=... scripts/install-mcp.sh

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN_DIR_PREFIX="$REPO_ROOT/mcp"

# Verify all binaries exist
SERVERS=(proxmox terraform ansible kubectl cloudflare ssh)
missing=()
for srv in "${SERVERS[@]}"; do
  if [[ ! -x "$BIN_DIR_PREFIX/$srv/bin/$srv" ]]; then
    missing+=("$srv")
  fi
done
if [[ ${#missing[@]} -gt 0 ]]; then
  echo "ERROR: Missing binaries: ${missing[*]}"
  echo "Run 'make mcp-build' first."
  exit 1
fi

echo "Registering MCP servers with Claude Code..."

# proxmox — requires PROXMOX_HOST, PROXMOX_TOKEN_ID, PROXMOX_TOKEN_SECRET
claude mcp add \
  -e PROXMOX_HOST="${PROXMOX_HOST:-}" \
  -e PROXMOX_TOKEN_ID="${PROXMOX_TOKEN_ID:-}" \
  -e PROXMOX_TOKEN_SECRET="${PROXMOX_TOKEN_SECRET:-}" \
  proxmox -- "$BIN_DIR_PREFIX/proxmox/bin/proxmox"

# terraform — no credentials; reads provider config from terraform files
claude mcp add \
  terraform -- "$BIN_DIR_PREFIX/terraform/bin/terraform"

# ansible — inherits SSH keys and inventory from environment
claude mcp add \
  ansible -- "$BIN_DIR_PREFIX/ansible/bin/ansible"

# kubectl — KUBECONFIG optional; defaults to ~/.kube/config
claude mcp add \
  -e KUBECONFIG="${KUBECONFIG:-}" \
  kubectl -- "$BIN_DIR_PREFIX/kubectl/bin/kubectl"

# cloudflare — requires CF_API_TOKEN
claude mcp add \
  -e CF_API_TOKEN="${CF_API_TOKEN:-}" \
  cloudflare -- "$BIN_DIR_PREFIX/cloudflare/bin/cloudflare"

# ssh — no credentials; uses SSH agent or key files
claude mcp add \
  ssh -- "$BIN_DIR_PREFIX/ssh/bin/ssh"

echo ""
echo "Done. Registered servers:"
claude mcp list
