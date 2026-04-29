package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func args(req mcp.CallToolRequest) map[string]any {
	m, _ := req.Params.Arguments.(map[string]any)
	return m
}

func strArg(req mcp.CallToolRequest, key string) string {
	v, _ := args(req)[key].(string)
	return v
}

func errResult(err error) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("error: " + err.Error()), nil
}

func runCmd(name string, cmdArgs []string, dir string) (string, error) {
	cmd := exec.Command(name, cmdArgs...)
	cmd.Env = os.Environ()
	if dir != "" {
		cmd.Dir = dir
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// ansible_run_playbook — run a playbook against an inventory.
func handleRunPlaybook(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	playbook := strArg(req, "playbook")
	inventory := strArg(req, "inventory")
	extraVarsJSON := strArg(req, "extra_vars")
	tags := strArg(req, "tags")
	limit := strArg(req, "limit")

	cmdArgs := []string{playbook, "-i", inventory}

	if extraVarsJSON != "" {
		var m map[string]any
		if err := json.Unmarshal([]byte(extraVarsJSON), &m); err != nil {
			return errResult(fmt.Errorf("extra_vars must be valid JSON: %w", err))
		}
		cmdArgs = append(cmdArgs, "--extra-vars", extraVarsJSON)
	}
	if tags != "" {
		cmdArgs = append(cmdArgs, "--tags", tags)
	}
	if limit != "" {
		cmdArgs = append(cmdArgs, "--limit", limit)
	}

	out, err := runCmd("ansible-playbook", cmdArgs, "")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("ansible-playbook failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

// ansible_ping — test SSH connectivity to all hosts in the inventory.
func handlePing(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	inventory := strArg(req, "inventory")
	limit := strArg(req, "limit")

	cmdArgs := []string{"all", "-i", inventory, "-m", "ping"}
	if limit != "" {
		cmdArgs = append(cmdArgs, "--limit", limit)
	}

	out, err := runCmd("ansible", cmdArgs, "")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("ansible ping failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

// ansible_list_playbooks — list .yml playbooks in a directory.
func handleListPlaybooks(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dir := strArg(req, "directory")
	if dir == "" {
		dir = "."
	}

	matches, err := filepath.Glob(filepath.Join(dir, "*.yml"))
	if err != nil {
		return errResult(err)
	}
	matches2, _ := filepath.Glob(filepath.Join(dir, "**", "*.yml"))
	matches = append(matches, matches2...)

	if len(matches) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("no .yml files found in %s", dir)), nil
	}
	return mcp.NewToolResultText(strings.Join(matches, "\n")), nil
}

// ansible_inventory_list — show the parsed inventory as JSON.
func handleInventoryList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	inventory := strArg(req, "inventory")

	out, err := runCmd("ansible-inventory", []string{"-i", inventory, "--list"}, "")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("ansible-inventory failed:\n%s", out)), nil
	}

	var v any
	if json.Unmarshal([]byte(strings.TrimSpace(out)), &v) == nil {
		pretty, _ := json.MarshalIndent(v, "", "  ")
		return mcp.NewToolResultText(string(pretty)), nil
	}
	return mcp.NewToolResultText(out), nil
}

func main() {
	s := server.NewMCPServer("ansible-mcp", "0.1.0",
		server.WithToolCapabilities(true),
	)

	s.AddTool(
		mcp.NewTool("ansible_run_playbook",
			mcp.WithDescription("Run an Ansible playbook against an inventory. Credentials and SSH keys are inherited from the environment."),
			mcp.WithString("playbook", mcp.Required(), mcp.Description("Path to the playbook .yml file")),
			mcp.WithString("inventory", mcp.Required(), mcp.Description("Path to the inventory file or directory")),
			mcp.WithString("extra_vars", mcp.Description("JSON object of extra variables to pass via --extra-vars (optional)")),
			mcp.WithString("tags", mcp.Description("Comma-separated list of tags to run (optional)")),
			mcp.WithString("limit", mcp.Description("Limit execution to a host pattern, e.g. 'webservers' or '192.168.86.20' (optional)")),
		),
		handleRunPlaybook,
	)

	s.AddTool(
		mcp.NewTool("ansible_ping",
			mcp.WithDescription("Test SSH connectivity to hosts in the inventory using the Ansible ping module."),
			mcp.WithString("inventory", mcp.Required(), mcp.Description("Path to the inventory file or directory")),
			mcp.WithString("limit", mcp.Description("Limit to a host pattern (optional)")),
		),
		handlePing,
	)

	s.AddTool(
		mcp.NewTool("ansible_list_playbooks",
			mcp.WithDescription("List available Ansible playbooks in a directory."),
			mcp.WithString("directory", mcp.Description("Directory to search. Defaults to current directory.")),
		),
		handleListPlaybooks,
	)

	s.AddTool(
		mcp.NewTool("ansible_inventory_list",
			mcp.WithDescription("Parse and display the inventory as structured JSON, showing all hosts and groups."),
			mcp.WithString("inventory", mcp.Required(), mcp.Description("Path to the inventory file or directory")),
		),
		handleInventoryList,
	)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
