package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/crypto/ssh"
)

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

func handleExec(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	host := strArg(req, "host")
	user := strArg(req, "user")
	command := strArg(req, "command")
	keyPath := strArg(req, "private_key_path")
	port := intArg(req, "port")
	timeoutSec := intArg(req, "timeout")

	if port == 0 {
		port = 22
	}
	if timeoutSec == 0 {
		timeoutSec = 30
	}
	if keyPath == "" {
		keyPath = filepath.Join(os.Getenv("HOME"), ".ssh", "id_ed25519")
		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			keyPath = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
		}
	}

	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return errResult(fmt.Errorf("read private key %s: %w", keyPath, err))
	}

	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return errResult(fmt.Errorf("parse private key: %w", err))
	}

	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		// Self-signed host keys are the norm in homelabs; skip verification.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec
		Timeout:         time.Duration(timeoutSec) * time.Second,
	}

	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return errResult(fmt.Errorf("connect to %s: %w", addr, err))
	}
	defer conn.Close()

	sess, err := conn.NewSession()
	if err != nil {
		return errResult(fmt.Errorf("open session: %w", err))
	}
	defer sess.Close()

	var stdout, stderr bytes.Buffer
	sess.Stdout = &stdout
	sess.Stderr = &stderr

	runErr := sess.Run(command)

	var out strings.Builder
	if stdout.Len() > 0 {
		fmt.Fprintf(&out, "stdout:\n%s\n", stdout.String())
	}
	if stderr.Len() > 0 {
		fmt.Fprintf(&out, "stderr:\n%s\n", stderr.String())
	}
	if runErr != nil {
		fmt.Fprintf(&out, "exit error: %v", runErr)
	} else {
		fmt.Fprintf(&out, "exit: 0")
	}

	return mcp.NewToolResultText(out.String()), nil
}

func errResult(err error) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("error: " + err.Error()), nil
}

func main() {
	s := server.NewMCPServer("ssh-mcp", "0.1.0",
		server.WithToolCapabilities(true),
	)

	s.AddTool(
		mcp.NewTool("ssh_exec",
			mcp.WithDescription("Execute a shell command on a remote host over SSH and return stdout, stderr, and exit status. Uses key-based authentication."),
			mcp.WithString("host", mcp.Required(), mcp.Description("Hostname or IP address of the remote machine")),
			mcp.WithString("user", mcp.Required(), mcp.Description("SSH user to connect as (e.g. root)")),
			mcp.WithString("command", mcp.Required(), mcp.Description("Shell command to execute")),
			mcp.WithString("private_key_path", mcp.Description("Path to SSH private key. Defaults to ~/.ssh/id_ed25519 then ~/.ssh/id_rsa")),
			mcp.WithNumber("port", mcp.Description("SSH port. Default: 22")),
			mcp.WithNumber("timeout", mcp.Description("Command timeout in seconds. Default: 30")),
		),
		handleExec,
	)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
