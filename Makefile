.PHONY: mcp-build mcp-install mcp-clean website-build website-dev website-deploy help

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
	@echo "  website-dev    Run Hugo dev server"
	@echo "  website-deploy Build and rsync to monitoring LXC (requires DEPLOY_HOST)"

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

HUGO := $(shell command -v hugo 2>/dev/null || echo /tmp/hugo)

website-build:
	cd website && $(HUGO) --minify

website-dev:
	cd website && $(HUGO) server --buildDrafts

# DEPLOY_HOST: SSH target for the monitoring LXC (e.g. root@192.168.86.25)
DEPLOY_HOST ?= root@192.168.86.25
DEPLOY_PATH := /opt/monitoring/homelab-site/build/html

website-deploy: website-build
	rsync -az --delete website/public/ $(DEPLOY_HOST):$(DEPLOY_PATH)/
	ssh $(DEPLOY_HOST) "docker compose -f /opt/monitoring/docker-compose.yml up -d --build homelab-site"
