package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

func newClient(host, apiToken string, skipTLS bool) *client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLS}, //nolint:gosec // self-signed certs are common in homelabs
	}
	return &client{
		baseURL:    fmt.Sprintf("https://%s:8006/api2/json", host),
		apiToken:   apiToken,
		httpClient: &http.Client{Transport: transport},
	}
}

type apiResponse struct {
	Data json.RawMessage `json:"data"`
}

func (c *client) checkCreds() error {
	if c.apiToken == "" || c.baseURL == "https://:8006/api2/json" {
		return fmt.Errorf("PROXMOX_HOST and PROXMOX_API_TOKEN env vars must be set")
	}
	return nil
}

func (c *client) get(path string) (json.RawMessage, error) {
	if err := c.checkCreds(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "PVEAPIToken="+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("proxmox API %d: %s", resp.StatusCode, string(body))
	}

	var result apiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result.Data, nil
}

func (c *client) post(path string, params map[string]string) (json.RawMessage, error) {
	if err := c.checkCreds(); err != nil {
		return nil, err
	}
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	req, err := http.NewRequest("POST", c.baseURL+path, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "PVEAPIToken="+c.apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("proxmox API %d: %s", resp.StatusCode, string(body))
	}

	var result apiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result.Data, nil
}

func (c *client) delete(path string) error {
	if err := c.checkCreds(); err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", c.baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "PVEAPIToken="+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("proxmox API %d: %s", resp.StatusCode, string(body))
	}
	return nil
}
