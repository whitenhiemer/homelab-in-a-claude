---
title: "How the MCP Servers Work"
date: 2026-03-25
---

The six MCP servers in this project are what give Claude Code direct access to your infrastructure. This post explains how they work, why they're structured the way they are, and how to add your own.

---

## What MCP is

MCP (Model Context Protocol) is a protocol for giving language models access to external tools. In Claude Code, you register MCP servers with `claude mcp add`. When Claude wants to use a tool, it sends a JSON request to the server over stdin; the server responds with the result on stdout.

Every MCP server in this project is a Go binary that:
1. Registers a set of named tools with descriptions and parameter schemas
2. Reads requests from stdin
3. Executes the underlying operation (API call, shell command, etc.)
4. Returns results as structured JSON

Claude never touches your infrastructure directly — it always goes through one of these servers, which means you can audit exactly what operations are available.

---

## The protocol constraint that shapes everything

MCP requires a handshake before any tools can be called. The client sends `initialize`, the server responds with its capabilities, and only then can tools be invoked.

This means a server that calls `os.Exit(1)` at startup — because an env var isn't set, for example — breaks the protocol. Claude Code sees the server die before the handshake and reports a connection error instead of a useful "PROXMOX_HOST is not set" message.

The correct pattern is: start the server unconditionally, complete the handshake, and only check credentials when a tool is actually called. Every server in this project defers credential validation to the first tool invocation and returns the error as a proper tool result with `isError: true`.

```go
// Wrong — kills the process before initialize handshake
func main() {
    if os.Getenv("PROXMOX_API_TOKEN") == "" {
        log.Fatal("PROXMOX_API_TOKEN must be set")
    }
    // ...
}

// Right — check at call time, return structured error
func (c *client) checkCreds() error {
    if c.apiToken == "" {
        return fmt.Errorf("PROXMOX_API_TOKEN env var must be set")
    }
    return nil
}
```

---

## The six servers

### proxmox

Wraps the Proxmox VE REST API. Tools:
- `list_nodes` — list cluster nodes and their status
- `list_vms` — list VMs across all nodes
- `list_containers` — list LXC containers
- `get_vm` / `get_container` — get config for a specific VM or LXC
- `start_vm` / `stop_vm` / `start_container` / `stop_container` — power management
- `exec_container` — run a command in a container via the Proxmox API

Credentials: `PROXMOX_HOST`, `PROXMOX_API_TOKEN`. Set `PROXMOX_INSECURE=true` if you're using a self-signed cert (common with Proxmox default installs).

### cloudflare

Wraps the Cloudflare v4 API. Tools:
- `list_zones` — list DNS zones
- `list_dns_records` — list records in a zone
- `create_dns_record` / `update_dns_record` / `delete_dns_record`
- `get_zone_id` — resolve a zone name to its ID

Credential: `CF_API_TOKEN`.

### ssh

Opens SSH connections to hosts and runs commands. Tools:
- `run_command` — run a shell command on a remote host
- `read_file` / `write_file` — read or write a file over SSH
- `upload_file` — upload a local file to a remote path

Auth: tries `SSH_AUTH_SOCK` (SSH agent) first, then searches `~/.ssh/id_ed25519`, `~/.ssh/id_ansible`, `~/.ssh/id_rsa` in order.

The SSH agent path is the most reliable in practice — if you've already authenticated in your shell session, the server piggybacks on that.

### terraform

Runs Terraform commands in a given directory. Tools:
- `init` — terraform init
- `plan` — terraform plan, returns the plan output
- `apply` — terraform apply (requires explicit confirmation flag)
- `destroy` — terraform destroy (requires explicit confirmation flag)
- `output` — get outputs from the state
- `show` — show current state

The `apply` and `destroy` tools require `confirmed: true` in the parameters — Claude must explicitly pass this flag, which forces it to ask you before running destructive operations.

### ansible

Runs Ansible playbooks. Tools:
- `run_playbook` — run a playbook with optional extra vars and tags
- `run_adhoc` — run an ad-hoc module on a host pattern
- `list_inventory` — list hosts from an inventory file

### kubectl

Runs kubectl commands against a kubeconfig. Tools:
- `get` — kubectl get with optional namespace and label selectors
- `describe` — kubectl describe
- `apply` — apply a manifest from a string or file
- `delete` — delete a resource
- `logs` — get pod logs
- `exec` — exec into a pod

Credential: `KUBECONFIG` env var pointing to your kubeconfig file.

---

## Adding a new server

The pattern is consistent across all six. Here's the skeleton:

```go
package main

import (
    "context"
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

type srv struct{ /* your client */ }

func main() {
    s := &srv{}
    mcpServer := server.NewMCPServer("my-tool", "1.0.0")

    mcpServer.AddTool(
        mcp.NewTool("do_thing",
            mcp.WithDescription("Does a thing"),
            mcp.WithString("param", mcp.Required(), mcp.Description("The param")),
        ),
        s.doThing,
    )

    server.ServeStdio(mcpServer)
}

func (s *srv) doThing(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    param := req.Params.Arguments["param"].(string)
    // do the thing
    return mcp.NewToolResultText("done: " + param), nil
}

func errResult(err error) (*mcp.CallToolResult, error) {
    return mcp.NewToolResultError(err.Error()), nil
}
```

Key points:
- Always use `mcp.NewToolResultError()` for errors, not `mcp.NewToolResultText()`. The `isError: true` flag is how Claude knows the operation failed.
- Don't call `os.Exit` before `server.ServeStdio()`.
- Keep tool names lowercase with underscores — Claude handles them more reliably than camelCase.
- Write descriptions as imperative sentences: "List all VMs on the cluster", not "Lists VMs" or "VM listing tool". Claude uses these to decide which tool to call.

Add it to `mcp/` and wire it into `Makefile` and `scripts/install-mcp.sh`.

---

## Debugging a server

If a tool isn't working, the fastest path is to run the server directly and send it a test request:

```bash
# Build and run the server
go build -o /tmp/test-server ./mcp/proxmox
PROXMOX_HOST=192.168.1.10 PROXMOX_API_TOKEN="root@pam!claude=..." /tmp/test-server

# In another terminal, send a raw MCP request
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | /tmp/test-server
```

Claude Code also logs MCP traffic if you run it with `--debug`. The logs show exactly what JSON is being sent and received for each tool call.
