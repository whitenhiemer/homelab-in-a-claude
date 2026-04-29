package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const cfAPI = "https://api.cloudflare.com/client/v4"

type cfClient struct {
	token      string
	httpClient *http.Client
}

func newCFClient(token string) *cfClient {
	return &cfClient{
		token:      token,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *cfClient) do(method, path string, body string) (json.RawMessage, error) {
	if c.token == "" {
		return nil, fmt.Errorf("CF_API_TOKEN env var must be set")
	}
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, cfAPI+path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Success bool              `json:"success"`
		Errors  []map[string]any  `json:"errors"`
		Result  json.RawMessage   `json:"result"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if !result.Success {
		errs, _ := json.Marshal(result.Errors)
		return nil, fmt.Errorf("cloudflare API error: %s", string(errs))
	}
	return result.Result, nil
}

func (c *cfClient) get(path string) (json.RawMessage, error) {
	return c.do("GET", path, "")
}

func (c *cfClient) post(path string, payload any) (json.RawMessage, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return c.do("POST", path, string(b))
}

func (c *cfClient) patch(path string, payload any) (json.RawMessage, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return c.do("PATCH", path, string(b))
}

func (c *cfClient) delete(path string) error {
	_, err := c.do("DELETE", path, "")
	return err
}

// ---- helpers ----

type srv struct{ cf *cfClient }

func args(req mcp.CallToolRequest) map[string]any {
	m, _ := req.Params.Arguments.(map[string]any)
	return m
}

func strArg(req mcp.CallToolRequest, key string) string {
	v, _ := args(req)[key].(string)
	return v
}

func boolArg(req mcp.CallToolRequest, key string) bool {
	v, _ := args(req)[key].(bool)
	return v
}

func resultJSON(v any) (*mcp.CallToolResult, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(b)), nil
}

func errResult(err error) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError(err.Error()), nil
}

// ---- tool handlers ----

func (s *srv) handleListZones(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	data, err := s.cf.get("/zones?per_page=50")
	if err != nil {
		return errResult(err)
	}
	var zones []map[string]any
	json.Unmarshal(data, &zones) //nolint:errcheck
	// Trim to useful fields
	out := make([]map[string]any, 0, len(zones))
	for _, z := range zones {
		out = append(out, map[string]any{
			"id":     z["id"],
			"name":   z["name"],
			"status": z["status"],
		})
	}
	return resultJSON(out)
}

func (s *srv) handleListRecords(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	zoneID := strArg(req, "zone_id")
	recType := strArg(req, "type")

	path := fmt.Sprintf("/zones/%s/dns_records?per_page=100", zoneID)
	if recType != "" {
		path += "&type=" + recType
	}

	data, err := s.cf.get(path)
	if err != nil {
		return errResult(err)
	}
	var records []map[string]any
	json.Unmarshal(data, &records) //nolint:errcheck
	out := make([]map[string]any, 0, len(records))
	for _, r := range records {
		out = append(out, map[string]any{
			"id":      r["id"],
			"type":    r["type"],
			"name":    r["name"],
			"content": r["content"],
			"proxied": r["proxied"],
			"ttl":     r["ttl"],
		})
	}
	return resultJSON(out)
}

func (s *srv) handleCreateRecord(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	zoneID := strArg(req, "zone_id")

	payload := map[string]any{
		"type":    strArg(req, "type"),
		"name":    strArg(req, "name"),
		"content": strArg(req, "content"),
		"proxied": boolArg(req, "proxied"),
		"ttl":     1, // auto
	}

	data, err := s.cf.post("/zones/"+zoneID+"/dns_records", payload)
	if err != nil {
		return errResult(err)
	}
	var record map[string]any
	json.Unmarshal(data, &record) //nolint:errcheck
	return mcp.NewToolResultText(fmt.Sprintf("Created record id=%s: %s %s → %s",
		record["id"], record["type"], record["name"], record["content"])), nil
}

func (s *srv) handleUpdateRecord(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	zoneID := strArg(req, "zone_id")
	recordID := strArg(req, "record_id")

	payload := map[string]any{
		"content": strArg(req, "content"),
		"proxied": boolArg(req, "proxied"),
	}

	data, err := s.cf.patch("/zones/"+zoneID+"/dns_records/"+recordID, payload)
	if err != nil {
		return errResult(err)
	}
	var record map[string]any
	json.Unmarshal(data, &record) //nolint:errcheck
	return mcp.NewToolResultText(fmt.Sprintf("Updated record %s → %s", record["name"], record["content"])), nil
}

func (s *srv) handleDeleteRecord(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	zoneID := strArg(req, "zone_id")
	recordID := strArg(req, "record_id")

	if err := s.cf.delete("/zones/" + zoneID + "/dns_records/" + recordID); err != nil {
		return errResult(err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("Deleted record %s from zone %s", recordID, zoneID)), nil
}

func main() {
	s := &srv{cf: newCFClient(os.Getenv("CF_API_TOKEN"))}
	mcpServer := server.NewMCPServer("cloudflare-mcp", "0.1.0",
		server.WithToolCapabilities(true),
	)

	mcpServer.AddTool(
		mcp.NewTool("cloudflare_list_zones",
			mcp.WithDescription("List all Cloudflare zones (domains) in the account."),
		),
		s.handleListZones,
	)

	mcpServer.AddTool(
		mcp.NewTool("cloudflare_list_records",
			mcp.WithDescription("List DNS records for a Cloudflare zone."),
			mcp.WithString("zone_id", mcp.Required(), mcp.Description("Zone ID from cloudflare_list_zones")),
			mcp.WithString("type", mcp.Description("Filter by record type: A, AAAA, CNAME, TXT, MX, etc. (optional)")),
		),
		s.handleListRecords,
	)

	mcpServer.AddTool(
		mcp.NewTool("cloudflare_create_record",
			mcp.WithDescription("Create a DNS record in a Cloudflare zone."),
			mcp.WithString("zone_id", mcp.Required(), mcp.Description("Zone ID")),
			mcp.WithString("type", mcp.Required(), mcp.Description("Record type: A, AAAA, CNAME, TXT, MX")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Record name, e.g. homelab.woodhead.tech or @ for apex")),
			mcp.WithString("content", mcp.Required(), mcp.Description("Record value, e.g. an IP address or target hostname")),
			mcp.WithBoolean("proxied", mcp.Description("Whether to proxy traffic through Cloudflare (orange cloud). Default false.")),
		),
		s.handleCreateRecord,
	)

	mcpServer.AddTool(
		mcp.NewTool("cloudflare_update_record",
			mcp.WithDescription("Update the content or proxy status of an existing DNS record."),
			mcp.WithString("zone_id", mcp.Required(), mcp.Description("Zone ID")),
			mcp.WithString("record_id", mcp.Required(), mcp.Description("Record ID from cloudflare_list_records")),
			mcp.WithString("content", mcp.Required(), mcp.Description("New record value")),
			mcp.WithBoolean("proxied", mcp.Description("Whether to proxy through Cloudflare")),
		),
		s.handleUpdateRecord,
	)

	mcpServer.AddTool(
		mcp.NewTool("cloudflare_delete_record",
			mcp.WithDescription("Delete a DNS record from a Cloudflare zone."),
			mcp.WithString("zone_id", mcp.Required(), mcp.Description("Zone ID")),
			mcp.WithString("record_id", mcp.Required(), mcp.Description("Record ID from cloudflare_list_records")),
		),
		s.handleDeleteRecord,
	)

	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
