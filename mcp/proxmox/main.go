package main

import (
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type srv struct {
	pve *client
}

func main() {
	s := &srv{pve: newClient(
		os.Getenv("PROXMOX_HOST"),
		os.Getenv("PROXMOX_API_TOKEN"),
		os.Getenv("PROXMOX_INSECURE") == "true",
	)}

	mcpServer := server.NewMCPServer("proxmox-mcp", "0.1.0",
		server.WithToolCapabilities(true),
	)

	mcpServer.AddTool(
		mcp.NewTool("proxmox_list_nodes",
			mcp.WithDescription("List all nodes in the Proxmox cluster with CPU, memory, and disk usage."),
		),
		s.handleListNodes,
	)

	mcpServer.AddTool(
		mcp.NewTool("proxmox_list_vms",
			mcp.WithDescription("List all VMs and LXC containers on a Proxmox node."),
			mcp.WithString("node", mcp.Required(), mcp.Description("Proxmox node name (e.g. pve1)")),
		),
		s.handleListVMs,
	)

	mcpServer.AddTool(
		mcp.NewTool("proxmox_get_vm_status",
			mcp.WithDescription("Get current status and resource usage for a specific VM or LXC container."),
			mcp.WithString("node", mcp.Required(), mcp.Description("Proxmox node name")),
			mcp.WithNumber("vmid", mcp.Required(), mcp.Description("VM or container ID")),
			mcp.WithString("type", mcp.Required(), mcp.Description("'qemu' for VMs, 'lxc' for containers")),
		),
		s.handleGetVMStatus,
	)

	mcpServer.AddTool(
		mcp.NewTool("proxmox_create_lxc",
			mcp.WithDescription("Create a new LXC container. Use proxmox_list_content to find available templates first."),
			mcp.WithString("node", mcp.Required(), mcp.Description("Proxmox node name")),
			mcp.WithNumber("vmid", mcp.Required(), mcp.Description("Container ID (100–999999)")),
			mcp.WithString("hostname", mcp.Required(), mcp.Description("Container hostname")),
			mcp.WithString("ostemplate", mcp.Required(), mcp.Description("Template path, e.g. local:vztmpl/debian-12-standard_12.7-1_amd64.tar.zst")),
			mcp.WithString("memory", mcp.Required(), mcp.Description("RAM in MB, e.g. 2048")),
			mcp.WithString("cores", mcp.Required(), mcp.Description("Number of CPU cores")),
			mcp.WithString("rootfs", mcp.Required(), mcp.Description("Root disk, e.g. local-lvm:8 for 8GB on local-lvm")),
			mcp.WithString("net0", mcp.Required(), mcp.Description("Network config, e.g. name=eth0,bridge=vmbr0,ip=192.168.86.50/24,gw=192.168.86.1")),
			mcp.WithString("password", mcp.Required(), mcp.Description("Root password for the container")),
			mcp.WithString("ssh_public_keys", mcp.Description("SSH public key to inject (optional but recommended)")),
			mcp.WithString("nameserver", mcp.Description("DNS server IP (optional)")),
		),
		s.handleCreateLXC,
	)

	mcpServer.AddTool(
		mcp.NewTool("proxmox_create_vm",
			mcp.WithDescription("Create a new QEMU VM. Use proxmox_list_content to find available ISOs first."),
			mcp.WithString("node", mcp.Required(), mcp.Description("Proxmox node name")),
			mcp.WithNumber("vmid", mcp.Required(), mcp.Description("VM ID (100–999999)")),
			mcp.WithString("name", mcp.Required(), mcp.Description("VM name")),
			mcp.WithString("memory", mcp.Required(), mcp.Description("RAM in MB")),
			mcp.WithString("cores", mcp.Required(), mcp.Description("Number of CPU cores")),
			mcp.WithString("disk", mcp.Required(), mcp.Description("Primary disk, e.g. local-lvm:50 for 50GB")),
			mcp.WithString("net0", mcp.Required(), mcp.Description("Network config, e.g. virtio,bridge=vmbr0")),
			mcp.WithString("cdrom", mcp.Description("ISO path for installation, e.g. local:iso/debian-12.iso")),
			mcp.WithString("balloon", mcp.Description("Memory balloon minimum in MB (optional)")),
		),
		s.handleCreateVM,
	)

	mcpServer.AddTool(
		mcp.NewTool("proxmox_start",
			mcp.WithDescription("Start a stopped VM or LXC container."),
			mcp.WithString("node", mcp.Required(), mcp.Description("Proxmox node name")),
			mcp.WithNumber("vmid", mcp.Required(), mcp.Description("VM or container ID")),
			mcp.WithString("type", mcp.Required(), mcp.Description("'qemu' or 'lxc'")),
		),
		s.handleStart,
	)

	mcpServer.AddTool(
		mcp.NewTool("proxmox_stop",
			mcp.WithDescription("Stop a running VM or LXC container (immediate power-off)."),
			mcp.WithString("node", mcp.Required(), mcp.Description("Proxmox node name")),
			mcp.WithNumber("vmid", mcp.Required(), mcp.Description("VM or container ID")),
			mcp.WithString("type", mcp.Required(), mcp.Description("'qemu' or 'lxc'")),
		),
		s.handleStop,
	)

	mcpServer.AddTool(
		mcp.NewTool("proxmox_delete",
			mcp.WithDescription("Delete a VM or LXC container. The guest must be stopped first."),
			mcp.WithString("node", mcp.Required(), mcp.Description("Proxmox node name")),
			mcp.WithNumber("vmid", mcp.Required(), mcp.Description("VM or container ID")),
			mcp.WithString("type", mcp.Required(), mcp.Description("'qemu' or 'lxc'")),
		),
		s.handleDelete,
	)

	mcpServer.AddTool(
		mcp.NewTool("proxmox_list_storage",
			mcp.WithDescription("List storage pools on a node with available space and supported content types."),
			mcp.WithString("node", mcp.Required(), mcp.Description("Proxmox node name")),
		),
		s.handleListStorage,
	)

	mcpServer.AddTool(
		mcp.NewTool("proxmox_list_content",
			mcp.WithDescription("List available ISOs or container templates in a storage pool."),
			mcp.WithString("node", mcp.Required(), mcp.Description("Proxmox node name")),
			mcp.WithString("storage", mcp.Required(), mcp.Description("Storage pool name, e.g. local")),
			mcp.WithString("content_type", mcp.Required(), mcp.Description("'iso' for ISO images, 'vztmpl' for LXC templates")),
		),
		s.handleListContent,
	)

	mcpServer.AddTool(
		mcp.NewTool("proxmox_get_task",
			mcp.WithDescription("Check the status of an async task (UPID) returned by create, start, or stop operations."),
			mcp.WithString("node", mcp.Required(), mcp.Description("Proxmox node name")),
			mcp.WithString("upid", mcp.Required(), mcp.Description("Task UPID returned by the triggering operation")),
		),
		s.handleGetTask,
	)

	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
