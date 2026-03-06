package vm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

// cmdOutput holds the result of a command execution.
type cmdOutput struct {
	Stdout string
	Stderr string
}

// isLocal returns true when server points to the same host this process runs on.
func isLocal(server *types.Server) bool {
	if server.IPV4 == "" {
		return true
	}
	return server.IPV4 == "127.0.0.1" || server.IPV4 == "localhost" || server.IPV4 == ""
}

// runCmd executes a shell command on the server. If the server is local the
// command runs via os/exec; otherwise it is dispatched over SSH.
func (m *VmManager) runCmd(ctx context.Context, server *types.Server, command string) (*cmdOutput, error) {
	if isLocal(server) {
		return m.runLocal(ctx, command)
	}
	return m.runSSH(ctx, server, command)
}

// runLocal executes command on the local host.
func (m *VmManager) runLocal(ctx context.Context, command string) (*cmdOutput, error) {
	log.Ctx(ctx).Debug().Str("cmd", command).Msg("vm: exec local")

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return &cmdOutput{Stdout: stdout.String(), Stderr: stderr.String()},
			fmt.Errorf("local command failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return &cmdOutput{Stdout: stdout.String(), Stderr: stderr.String()}, nil
}

// runSSH executes command on a remote server over SSH.
func (m *VmManager) runSSH(ctx context.Context, server *types.Server, command string) (*cmdOutput, error) {
	log.Ctx(ctx).Debug().Str("cmd", command).Str("host", server.IPV4).Msg("vm: exec ssh")

	client, err := m.sshClient(server)
	if err != nil {
		return nil, fmt.Errorf("ssh connect: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("ssh session: %w", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	done := make(chan error, 1)
	go func() { done <- session.Run(command) }()

	select {
	case <-ctx.Done():
		_ = session.Signal(ssh.SIGTERM)
		return nil, ctx.Err()
	case err := <-done:
		if err != nil {
			return &cmdOutput{Stdout: stdout.String(), Stderr: stderr.String()},
				fmt.Errorf("ssh command failed: %w: %s", err, strings.TrimSpace(stderr.String()))
		}
	}

	return &cmdOutput{Stdout: stdout.String(), Stderr: stderr.String()}, nil
}

// runSSHStream executes a command over SSH and streams stdout line-by-line
// through the returned io.ReadCloser. Caller must close when done.
func (m *VmManager) runSSHStream(ctx context.Context, server *types.Server, command string) (io.ReadCloser, func(), error) {
	client, err := m.sshClient(server)
	if err != nil {
		return nil, nil, fmt.Errorf("ssh connect: %w", err)
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("ssh session: %w", err)
	}

	pipe, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		client.Close()
		return nil, nil, fmt.Errorf("ssh stdout pipe: %w", err)
	}

	if err := session.Start(command); err != nil {
		session.Close()
		client.Close()
		return nil, nil, fmt.Errorf("ssh start: %w", err)
	}

	cleanup := func() {
		_ = session.Close()
		_ = client.Close()
	}

	return io.NopCloser(pipe), cleanup, nil
}

const sshKeyPath = "/data/cloudness/ssh/id_rsa"

// sshClient creates a new SSH client connection to the server.
func (m *VmManager) sshClient(server *types.Server) (*ssh.Client, error) {
	user := server.User
	if user == "" {
		user = "root"
	}

	port := server.Port
	if port == 0 {
		port = 22
	}

	keyBytes, err := os.ReadFile(sshKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key %s: %w", sshKeyPath, err)
	}

	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := net.JoinHostPort(server.IPV4, fmt.Sprintf("%d", port))
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", addr, err)
	}

	return client, nil
}
