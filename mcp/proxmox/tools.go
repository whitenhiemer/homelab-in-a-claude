package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func resultJSON(v any) (*mcp.CallToolResult, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(b)), nil
}

func resultError(err error) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError(err.Error()), nil
}

func args(req mcp.CallToolRequest) map[string]any {
	m, _ := req.Params.Arguments.(map[string]any)
	return m
}

func strArg(req mcp.CallToolRequest, key string) string {
	v, _ := args(req)[key].(string)
	return v
}

func intArg(req mcp.CallToolRequest, key string) int {
	switch v := args(req)[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	}
	return 0
}

// proxmox_list_nodes — list all nodes in the cluster with resource usage.
func (s *srv) handleListNodes(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	data, err := s.pve.get("/nodes")
	if err != nil {
		return resultError(err)
	}
	var nodes []map[string]any
	if err := json.Unmarshal(data, &nodes); err != nil {
		return resultError(err)
	}
	return resultJSON(nodes)
}

// proxmox_list_vms — list VMs and LXC containers on a node.
func (s *srv) handleListVMs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	node := strArg(req, "node")

	qemuData, err := s.pve.get("/nodes/" + node + "/qemu")
	if err != nil {
		return resultError(fmt.Errorf("list qemu: %w", err))
	}
	var qemuList []map[string]any
	json.Unmarshal(qemuData, &qemuList) //nolint:errcheck
	for _, vm := range qemuList {
		vm["type"] = "qemu"
	}

	lxcData, err := s.pve.get("/nodes/" + node + "/lxc")
	if err != nil {
		return resultError(fmt.Errorf("list lxc: %w", err))
	}
	var lxcList []map[string]any
	json.Unmarshal(lxcData, &lxcList) //nolint:errcheck
	for _, ct := range lxcList {
		ct["type"] = "lxc"
	}

	all := append(qemuList, lxcList...)
	return resultJSON(all)
}

// proxmox_get_vm_status — get current status and resource usage for a VM or LXC.
func (s *srv) handleGetVMStatus(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	node := strArg(req, "node")
	vmid := intArg(req, "vmid")
	vmtype := strArg(req, "type")

	path := fmt.Sprintf("/nodes/%s/%s/%d/status/current", node, vmtype, vmid)
	data, err := s.pve.get(path)
	if err != nil {
		return resultError(err)
	}
	var status map[string]any
	json.Unmarshal(data, &status) //nolint:errcheck
	return resultJSON(status)
}

// proxmox_create_lxc — create an LXC container.
func (s *srv) handleCreateLXC(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	node := strArg(req, "node")
	vmid := intArg(req, "vmid")

	params := map[string]string{
		"vmid":         fmt.Sprintf("%d", vmid),
		"hostname":     strArg(req, "hostname"),
		"ostemplate":   strArg(req, "ostemplate"),
		"memory":       strArg(req, "memory"),
		"cores":        strArg(req, "cores"),
		"rootfs":       strArg(req, "rootfs"),
		"net0":         strArg(req, "net0"),
		"password":     strArg(req, "password"),
		"unprivileged": "1",
		"start":        "0",
	}
	if ssh := strArg(req, "ssh_public_keys"); ssh != "" {
		params["ssh-public-keys"] = ssh
	}
	if ns := strArg(req, "nameserver"); ns != "" {
		params["nameserver"] = ns
	}

	data, err := s.pve.post("/nodes/"+node+"/lxc", params)
	if err != nil {
		return resultError(err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("LXC %d creation queued. Task: %s", vmid, string(data))), nil
}

// proxmox_create_vm — create a QEMU VM.
func (s *srv) handleCreateVM(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	node := strArg(req, "node")
	vmid := intArg(req, "vmid")

	params := map[string]string{
		"vmid":    fmt.Sprintf("%d", vmid),
		"name":    strArg(req, "name"),
		"memory":  strArg(req, "memory"),
		"cores":   strArg(req, "cores"),
		"sockets": "1",
		"cpu":     "x86-64-v2-AES",
		"bios":    "seabios",
		"scsi0":   strArg(req, "disk"),
		"net0":    strArg(req, "net0"),
		"boot":    "order=scsi0;ide2",
	}
	if iso := strArg(req, "cdrom"); iso != "" {
		params["ide2"] = iso + ",media=cdrom"
	}
	if balloon := strArg(req, "balloon"); balloon != "" {
		params["balloon"] = balloon
	}

	data, err := s.pve.post("/nodes/"+node+"/qemu", params)
	if err != nil {
		return resultError(err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("VM %d creation queued. Task: %s", vmid, string(data))), nil
}

// proxmox_start — start a VM or LXC.
func (s *srv) handleStart(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	node := strArg(req, "node")
	vmid := intArg(req, "vmid")
	vmtype := strArg(req, "type")

	path := fmt.Sprintf("/nodes/%s/%s/%d/status/start", node, vmtype, vmid)
	data, err := s.pve.post(path, nil)
	if err != nil {
		return resultError(err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("Start queued. Task: %s", string(data))), nil
}

// proxmox_stop — stop a VM or LXC.
func (s *srv) handleStop(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	node := strArg(req, "node")
	vmid := intArg(req, "vmid")
	vmtype := strArg(req, "type")

	path := fmt.Sprintf("/nodes/%s/%s/%d/status/stop", node, vmtype, vmid)
	data, err := s.pve.post(path, nil)
	if err != nil {
		return resultError(err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("Stop queued. Task: %s", string(data))), nil
}

// proxmox_delete — delete a VM or LXC (must be stopped first).
func (s *srv) handleDelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	node := strArg(req, "node")
	vmid := intArg(req, "vmid")
	vmtype := strArg(req, "type")

	path := fmt.Sprintf("/nodes/%s/%s/%d", node, vmtype, vmid)
	if err := s.pve.delete(path); err != nil {
		return resultError(err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("Deleted %s %d on %s", vmtype, vmid, node)), nil
}

// proxmox_list_storage — list storage pools on a node with usage stats.
func (s *srv) handleListStorage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	node := strArg(req, "node")
	data, err := s.pve.get("/nodes/" + node + "/storage")
	if err != nil {
		return resultError(err)
	}
	var pools []map[string]any
	json.Unmarshal(data, &pools) //nolint:errcheck
	return resultJSON(pools)
}

// proxmox_list_content — list ISOs or templates in a storage pool.
func (s *srv) handleListContent(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	node := strArg(req, "node")
	storage := strArg(req, "storage")
	contentType := strArg(req, "content_type")

	path := fmt.Sprintf("/nodes/%s/storage/%s/content?content=%s", node, storage, contentType)
	data, err := s.pve.get(path)
	if err != nil {
		return resultError(err)
	}
	var items []map[string]any
	json.Unmarshal(data, &items) //nolint:errcheck
	return resultJSON(items)
}

// proxmox_get_task — check the status of an async task returned by create/start/stop.
func (s *srv) handleGetTask(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	node := strArg(req, "node")
	upid := strArg(req, "upid")

	data, err := s.pve.get(fmt.Sprintf("/nodes/%s/tasks/%s/status", node, upid))
	if err != nil {
		return resultError(err)
	}
	var status map[string]any
	json.Unmarshal(data, &status) //nolint:errcheck
	return resultJSON(status)
}
