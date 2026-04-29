.PHONY: mcp-build mcp-install mcp-clean website-build website-dev help

MCP_SERVERS := proxmox terraform ansible kubectl cloudflare ssh

# Use GOROOT from environment if set, otherwise rely on PATH
GO := $(if $(shell command -v go 2>/dev/null),go,$(error "go not found in PATH"))

help:
	@echo "homelab-in-a-claude"
	@echo ""
	@echo "  mcp-build     Build all MCP server binaries"
	@echo "  mcp-install   Build and register MCP servers with Claude Code"
	@echo "  mcp-clean     Remove built binaries"
	@echo "  website-build Build the Hugo website"
	@echo "  website-dev   Run Hugo dev server"

mcp-build:
	@for srv in $(MCP_SERVERS); do \
		echo "Building mcp/$$srv..."; \
		mkdir -p mcp/$$srv/bin; \
		(cd mcp/$$srv && $(GO) build -o bin/$$srv .) || exit 1; \
	done

mcp-install: mcp-build
	@chmod +x scripts/install-mcp.sh
	@scripts/install-mcp.sh

mcp-clean:
	@for srv in $(MCP_SERVERS); do \
		rm -f mcp/$$srv/bin/$$srv; \
	done

website-build:
	cd website && hugo --minify

website-dev:
	cd website && hugo server --buildDrafts
