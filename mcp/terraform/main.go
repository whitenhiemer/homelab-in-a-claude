package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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

// runTerraform executes terraform with the given args in workDir, passing
// through the caller's environment so credentials set by the user are visible.
func runTerraform(workDir string, tfArgs []string, extraEnv []string) (string, error) {
	cmd := exec.Command("terraform", tfArgs...)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), extraEnv...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	return out.String(), err
}

// buildVarFlags turns a JSON-encoded vars string into -var=k=v flags.
func buildVarFlags(varsJSON string) ([]string, error) {
	if varsJSON == "" {
		return nil, nil
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(varsJSON), &m); err != nil {
		return nil, fmt.Errorf("vars must be a JSON object of string values: %w", err)
	}
	flags := make([]string, 0, len(m))
	for k, v := range m {
		flags = append(flags, fmt.Sprintf("-var=%s=%s", k, v))
	}
	return flags, nil
}

func handleInit(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dir := strArg(req, "working_dir")
	out, err := runTerraform(dir, []string{"init", "-no-color"}, nil)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("terraform init failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

func handlePlan(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dir := strArg(req, "working_dir")
	varFile := strArg(req, "var_file")
	varsJSON := strArg(req, "vars")

	tfArgs := []string{"plan", "-no-color"}
	if varFile != "" {
		tfArgs = append(tfArgs, "-var-file="+varFile)
	}
	varFlags, err := buildVarFlags(varsJSON)
	if err != nil {
		return errResult(err)
	}
	tfArgs = append(tfArgs, varFlags...)

	out, err := runTerraform(dir, tfArgs, nil)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("terraform plan failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

func handleApply(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dir := strArg(req, "working_dir")
	varFile := strArg(req, "var_file")
	varsJSON := strArg(req, "vars")

	tfArgs := []string{"apply", "-auto-approve", "-no-color"}
	if varFile != "" {
		tfArgs = append(tfArgs, "-var-file="+varFile)
	}
	varFlags, err := buildVarFlags(varsJSON)
	if err != nil {
		return errResult(err)
	}
	tfArgs = append(tfArgs, varFlags...)

	out, err := runTerraform(dir, tfArgs, nil)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("terraform apply failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

func handleDestroy(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dir := strArg(req, "working_dir")
	varFile := strArg(req, "var_file")

	tfArgs := []string{"destroy", "-auto-approve", "-no-color"}
	if varFile != "" {
		tfArgs = append(tfArgs, "-var-file="+varFile)
	}

	out, err := runTerraform(dir, tfArgs, nil)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("terraform destroy failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

func handleOutput(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dir := strArg(req, "working_dir")
	name := strArg(req, "output_name")

	tfArgs := []string{"output", "-json"}
	if name != "" {
		tfArgs = append(tfArgs, name)
	}

	out, err := runTerraform(dir, tfArgs, nil)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("terraform output failed:\n%s", out)), nil
	}

	// Pretty-print JSON output if possible
	var v any
	if json.Unmarshal([]byte(strings.TrimSpace(out)), &v) == nil {
		pretty, _ := json.MarshalIndent(v, "", "  ")
		return mcp.NewToolResultText(string(pretty)), nil
	}
	return mcp.NewToolResultText(out), nil
}

func handleShow(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dir := strArg(req, "working_dir")

	out, err := runTerraform(dir, []string{"show", "-no-color"}, nil)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("terraform show failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

func main() {
	s := server.NewMCPServer("terraform-mcp", "0.1.0",
		server.WithToolCapabilities(true),
	)

	s.AddTool(
		mcp.NewTool("terraform_init",
			mcp.WithDescription("Initialize a Terraform working directory. Run this before plan or apply in a recipe directory."),
			mcp.WithString("working_dir", mcp.Required(), mcp.Description("Path to the Terraform recipe directory")),
		),
		handleInit,
	)

	s.AddTool(
		mcp.NewTool("terraform_plan",
			mcp.WithDescription("Generate and show a Terraform execution plan. Always run this before apply so the user can review changes."),
			mcp.WithString("working_dir", mcp.Required(), mcp.Description("Path to the Terraform recipe directory")),
			mcp.WithString("var_file", mcp.Description("Path to a .tfvars file (optional)")),
			mcp.WithString("vars", mcp.Description("JSON object of variable overrides, e.g. {\"vmid\": \"150\", \"ip\": \"192.168.86.50\"} (optional)")),
		),
		handlePlan,
	)

	s.AddTool(
		mcp.NewTool("terraform_apply",
			mcp.WithDescription("Apply a Terraform plan. Only call this after showing the plan output to the user and getting confirmation."),
			mcp.WithString("working_dir", mcp.Required(), mcp.Description("Path to the Terraform recipe directory")),
			mcp.WithString("var_file", mcp.Description("Path to a .tfvars file (optional)")),
			mcp.WithString("vars", mcp.Description("JSON object of variable overrides (optional)")),
		),
		handleApply,
	)

	s.AddTool(
		mcp.NewTool("terraform_destroy",
			mcp.WithDescription("Destroy all resources managed by the Terraform state in this directory. Destructive — confirm with the user first."),
			mcp.WithString("working_dir", mcp.Required(), mcp.Description("Path to the Terraform recipe directory")),
			mcp.WithString("var_file", mcp.Description("Path to a .tfvars file (optional)")),
		),
		handleDestroy,
	)

	s.AddTool(
		mcp.NewTool("terraform_output",
			mcp.WithDescription("Read output values from the Terraform state. Useful for getting IPs or IDs of provisioned resources."),
			mcp.WithString("working_dir", mcp.Required(), mcp.Description("Path to the Terraform recipe directory")),
			mcp.WithString("output_name", mcp.Description("Specific output to retrieve. Omit to get all outputs.")),
		),
		handleOutput,
	)

	s.AddTool(
		mcp.NewTool("terraform_show",
			mcp.WithDescription("Show the current Terraform state in human-readable form."),
			mcp.WithString("working_dir", mcp.Required(), mcp.Description("Path to the Terraform recipe directory")),
		),
		handleShow,
	)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
