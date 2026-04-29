package main

import (
	"bytes"
	"context"
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

func runBin(bin string, binArgs []string, kubeconfig string) (string, error) {
	cmd := exec.Command(bin, binArgs...)
	cmd.Env = os.Environ()
	if kubeconfig != "" {
		cmd.Env = append(cmd.Env, "KUBECONFIG="+kubeconfig)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func kubeconfigArg(req mcp.CallToolRequest) string {
	return strArg(req, "kubeconfig")
}

// kubectl_get — get cluster resources.
func handleGet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	resource := strArg(req, "resource")
	name := strArg(req, "name")
	namespace := strArg(req, "namespace")

	cmdArgs := []string{"get", resource}
	if name != "" {
		cmdArgs = append(cmdArgs, name)
	}
	if namespace != "" {
		cmdArgs = append(cmdArgs, "-n", namespace)
	} else {
		cmdArgs = append(cmdArgs, "--all-namespaces")
	}
	cmdArgs = append(cmdArgs, "-o", "wide")

	out, err := runBin("kubectl", cmdArgs, kubeconfigArg(req))
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("kubectl get failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

// kubectl_apply — apply a manifest file or directory.
func handleApply(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := strArg(req, "path")
	namespace := strArg(req, "namespace")

	cmdArgs := []string{"apply", "-f", path}
	if namespace != "" {
		cmdArgs = append(cmdArgs, "-n", namespace)
	}

	out, err := runBin("kubectl", cmdArgs, kubeconfigArg(req))
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("kubectl apply failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

// kubectl_delete — delete a resource by type and name, or by manifest.
func handleDelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	resource := strArg(req, "resource")
	name := strArg(req, "name")
	namespace := strArg(req, "namespace")
	manifestPath := strArg(req, "manifest_path")

	var cmdArgs []string
	if manifestPath != "" {
		cmdArgs = []string{"delete", "-f", manifestPath}
	} else {
		cmdArgs = []string{"delete", resource, name}
	}
	if namespace != "" {
		cmdArgs = append(cmdArgs, "-n", namespace)
	}

	out, err := runBin("kubectl", cmdArgs, kubeconfigArg(req))
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("kubectl delete failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

// kubectl_describe — describe a resource in detail.
func handleDescribe(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	resource := strArg(req, "resource")
	name := strArg(req, "name")
	namespace := strArg(req, "namespace")

	cmdArgs := []string{"describe", resource}
	if name != "" {
		cmdArgs = append(cmdArgs, name)
	}
	if namespace != "" {
		cmdArgs = append(cmdArgs, "-n", namespace)
	}

	out, err := runBin("kubectl", cmdArgs, kubeconfigArg(req))
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("kubectl describe failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

// kubectl_logs — get logs from a pod.
func handleLogs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pod := strArg(req, "pod")
	namespace := strArg(req, "namespace")
	container := strArg(req, "container")
	tail := strArg(req, "tail")

	cmdArgs := []string{"logs", pod}
	if namespace != "" {
		cmdArgs = append(cmdArgs, "-n", namespace)
	}
	if container != "" {
		cmdArgs = append(cmdArgs, "-c", container)
	}
	if tail != "" {
		cmdArgs = append(cmdArgs, "--tail="+tail)
	} else {
		cmdArgs = append(cmdArgs, "--tail=100")
	}

	out, err := runBin("kubectl", cmdArgs, kubeconfigArg(req))
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("kubectl logs failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

// helm_install — install a Helm chart.
func handleHelmInstall(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	release := strArg(req, "release")
	chart := strArg(req, "chart")
	namespace := strArg(req, "namespace")
	valuesFile := strArg(req, "values_file")
	setValues := strArg(req, "set")

	cmdArgs := []string{"install", release, chart, "--create-namespace"}
	if namespace != "" {
		cmdArgs = append(cmdArgs, "-n", namespace)
	}
	if valuesFile != "" {
		cmdArgs = append(cmdArgs, "-f", valuesFile)
	}
	for _, kv := range strings.Split(setValues, ",") {
		kv = strings.TrimSpace(kv)
		if kv != "" {
			cmdArgs = append(cmdArgs, "--set", kv)
		}
	}

	out, err := runBin("helm", cmdArgs, kubeconfigArg(req))
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("helm install failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

// helm_upgrade — upgrade (or install) a Helm release.
func handleHelmUpgrade(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	release := strArg(req, "release")
	chart := strArg(req, "chart")
	namespace := strArg(req, "namespace")
	valuesFile := strArg(req, "values_file")
	setValues := strArg(req, "set")

	cmdArgs := []string{"upgrade", "--install", release, chart, "--create-namespace"}
	if namespace != "" {
		cmdArgs = append(cmdArgs, "-n", namespace)
	}
	if valuesFile != "" {
		cmdArgs = append(cmdArgs, "-f", valuesFile)
	}
	for _, kv := range strings.Split(setValues, ",") {
		kv = strings.TrimSpace(kv)
		if kv != "" {
			cmdArgs = append(cmdArgs, "--set", kv)
		}
	}

	out, err := runBin("helm", cmdArgs, kubeconfigArg(req))
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("helm upgrade failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

// helm_uninstall — remove a Helm release.
func handleHelmUninstall(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	release := strArg(req, "release")
	namespace := strArg(req, "namespace")

	cmdArgs := []string{"uninstall", release}
	if namespace != "" {
		cmdArgs = append(cmdArgs, "-n", namespace)
	}

	out, err := runBin("helm", cmdArgs, kubeconfigArg(req))
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("helm uninstall failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

// helm_list — list installed Helm releases.
func handleHelmList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace := strArg(req, "namespace")

	cmdArgs := []string{"list"}
	if namespace != "" {
		cmdArgs = append(cmdArgs, "-n", namespace)
	} else {
		cmdArgs = append(cmdArgs, "--all-namespaces")
	}

	out, err := runBin("helm", cmdArgs, kubeconfigArg(req))
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("helm list failed:\n%s", out)), nil
	}
	return mcp.NewToolResultText(out), nil
}

func main() {
	s := server.NewMCPServer("kubectl-mcp", "0.1.0",
		server.WithToolCapabilities(true),
	)

	kubeconfigDesc := mcp.Description("Path to kubeconfig file. Defaults to KUBECONFIG env var or ~/.kube/config.")

	s.AddTool(
		mcp.NewTool("kubectl_get",
			mcp.WithDescription("List Kubernetes resources of a given type. Returns wide output including node assignments and IPs."),
			mcp.WithString("resource", mcp.Required(), mcp.Description("Resource type, e.g. pods, nodes, services, deployments")),
			mcp.WithString("name", mcp.Description("Specific resource name (optional — omit to list all)")),
			mcp.WithString("namespace", mcp.Description("Namespace to query. Omit to query all namespaces.")),
			mcp.WithString("kubeconfig", kubeconfigDesc),
		),
		handleGet,
	)

	s.AddTool(
		mcp.NewTool("kubectl_apply",
			mcp.WithDescription("Apply a Kubernetes manifest file or directory to the cluster."),
			mcp.WithString("path", mcp.Required(), mcp.Description("Path to a manifest file or directory of manifests")),
			mcp.WithString("namespace", mcp.Description("Target namespace (optional)")),
			mcp.WithString("kubeconfig", kubeconfigDesc),
		),
		handleApply,
	)

	s.AddTool(
		mcp.NewTool("kubectl_delete",
			mcp.WithDescription("Delete a Kubernetes resource by type+name or by manifest file."),
			mcp.WithString("resource", mcp.Description("Resource type, e.g. pod, deployment (required if not using manifest_path)")),
			mcp.WithString("name", mcp.Description("Resource name (required if not using manifest_path)")),
			mcp.WithString("manifest_path", mcp.Description("Path to manifest file to delete by (alternative to resource+name)")),
			mcp.WithString("namespace", mcp.Description("Namespace (optional)")),
			mcp.WithString("kubeconfig", kubeconfigDesc),
		),
		handleDelete,
	)

	s.AddTool(
		mcp.NewTool("kubectl_describe",
			mcp.WithDescription("Show detailed information about a Kubernetes resource, including events."),
			mcp.WithString("resource", mcp.Required(), mcp.Description("Resource type, e.g. pod, node, service")),
			mcp.WithString("name", mcp.Description("Resource name (optional — omit to describe all of that type)")),
			mcp.WithString("namespace", mcp.Description("Namespace (optional)")),
			mcp.WithString("kubeconfig", kubeconfigDesc),
		),
		handleDescribe,
	)

	s.AddTool(
		mcp.NewTool("kubectl_logs",
			mcp.WithDescription("Retrieve logs from a pod. Returns the last 100 lines by default."),
			mcp.WithString("pod", mcp.Required(), mcp.Description("Pod name")),
			mcp.WithString("namespace", mcp.Description("Namespace (optional)")),
			mcp.WithString("container", mcp.Description("Container name for multi-container pods (optional)")),
			mcp.WithString("tail", mcp.Description("Number of lines from the end to return. Default: 100")),
			mcp.WithString("kubeconfig", kubeconfigDesc),
		),
		handleLogs,
	)

	s.AddTool(
		mcp.NewTool("helm_install",
			mcp.WithDescription("Install a Helm chart as a new release."),
			mcp.WithString("release", mcp.Required(), mcp.Description("Release name")),
			mcp.WithString("chart", mcp.Required(), mcp.Description("Chart reference, e.g. bitnami/nginx or ./charts/myapp")),
			mcp.WithString("namespace", mcp.Description("Target namespace (optional)")),
			mcp.WithString("values_file", mcp.Description("Path to a values.yaml override file (optional)")),
			mcp.WithString("set", mcp.Description("Comma-separated key=value pairs to pass via --set (optional)")),
			mcp.WithString("kubeconfig", kubeconfigDesc),
		),
		handleHelmInstall,
	)

	s.AddTool(
		mcp.NewTool("helm_upgrade",
			mcp.WithDescription("Upgrade an existing Helm release, or install it if it doesn't exist (--install)."),
			mcp.WithString("release", mcp.Required(), mcp.Description("Release name")),
			mcp.WithString("chart", mcp.Required(), mcp.Description("Chart reference")),
			mcp.WithString("namespace", mcp.Description("Target namespace (optional)")),
			mcp.WithString("values_file", mcp.Description("Path to a values.yaml override file (optional)")),
			mcp.WithString("set", mcp.Description("Comma-separated key=value pairs to pass via --set (optional)")),
			mcp.WithString("kubeconfig", kubeconfigDesc),
		),
		handleHelmUpgrade,
	)

	s.AddTool(
		mcp.NewTool("helm_uninstall",
			mcp.WithDescription("Remove a Helm release from the cluster."),
			mcp.WithString("release", mcp.Required(), mcp.Description("Release name to uninstall")),
			mcp.WithString("namespace", mcp.Description("Namespace the release is in (optional)")),
			mcp.WithString("kubeconfig", kubeconfigDesc),
		),
		handleHelmUninstall,
	)

	s.AddTool(
		mcp.NewTool("helm_list",
			mcp.WithDescription("List all installed Helm releases, optionally filtered by namespace."),
			mcp.WithString("namespace", mcp.Description("Namespace to filter by. Omit to list all namespaces.")),
			mcp.WithString("kubeconfig", kubeconfigDesc),
		),
		handleHelmList,
	)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
