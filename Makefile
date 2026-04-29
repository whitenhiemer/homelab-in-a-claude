.PHONY: mcp-build mcp-install mcp-clean website-build website-dev help

MCP_SERVERS := proxmox terraform ansible kubectl cloudflare ssh
CLAUDE_MCP_CONFIG := $(HOME)/.claude/mcp.json

help:
	@echo "homelab-in-a-claude"
	@echo ""
	@echo "  mcp-build     Build all MCP server binaries"
	@echo "  mcp-install   Register MCP servers with Claude Code"
	@echo "  mcp-clean     Remove built binaries"
	@echo "  website-build Build the Hugo website"
	@echo "  website-dev   Run Hugo dev server"

mcp-build:
	@for srv in $(MCP_SERVERS); do \
		echo "Building mcp/$$srv..."; \
		cd mcp/$$srv && go build -o bin/$$srv . && cd ../..; \
	done

mcp-install: mcp-build
	@echo "Registering MCP servers with Claude Code..."
	@node scripts/install-mcp.js

mcp-clean:
	@for srv in $(MCP_SERVERS); do \
		rm -f mcp/$$srv/bin/$$srv; \
	done

website-build:
	cd website && hugo --minify

website-dev:
	cd website && hugo server --buildDrafts
