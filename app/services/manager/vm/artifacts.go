package vm

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func (m *VmManager) ListArtifacts(ctx context.Context, server *types.Server, app *types.Application) ([]*types.Artifact, error) {
	identifier := app.GetIdentifierStr()
	filter := fmt.Sprintf("--filter label=%s --filter label=%s", managedByFilter(), appFilter(identifier))

	// List container names and IDs.
	cmd := fmt.Sprintf("docker ps -a --format '{{.ID}}|{{.Names}}' %s", filter)
	out, err := m.runCmd(ctx, server, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(out.Stdout), "\n")
	var artifacts []*types.Artifact

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		uid := parts[0]
		name := fmt.Sprintf("%s-%d", app.Name, i+1)
		if len(parts) == 2 && parts[1] != "" {
			name = parts[1]
		}
		artifacts = append(artifacts, &types.Artifact{
			UID:  uid,
			Name: name,
		})
	}

	return artifacts, nil
}

func (m *VmManager) TailLogs(ctx context.Context, server *types.Server, app *types.Application) (<-chan *types.ArtifactLogLine, <-chan error, error) {
	identifier := app.GetIdentifierStr()
	filter := fmt.Sprintf("--filter label=%s --filter label=%s", managedByFilter(), appFilter(identifier))

	// Get container IDs first.
	cmd := fmt.Sprintf("docker ps -aq %s", filter)
	out, err := m.runCmd(ctx, server, cmd)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list containers: %w", err)
	}

	containerIDs := strings.Fields(strings.TrimSpace(out.Stdout))
	if len(containerIDs) == 0 {
		logc := make(chan *types.ArtifactLogLine)
		errc := make(chan error)
		close(logc)
		close(errc)
		return logc, errc, nil
	}

	logc := make(chan *types.ArtifactLogLine)
	errc := make(chan error, 1)

	go func() {
		defer close(logc)
		for i, cid := range containerIDs {
			if err := m.streamContainerLogs(ctx, server, cid, fmt.Sprintf("%s-%d", app.Name, i), logc, errc); err != nil {
				log.Ctx(ctx).Error().Err(err).Str("container", cid).Msg("vm: error streaming logs")
				return
			}
		}
	}()

	return logc, errc, nil
}

func (m *VmManager) streamContainerLogs(ctx context.Context, server *types.Server, containerID, artifactUID string, logc chan<- *types.ArtifactLogLine, errc chan<- error) error {
	logCmd := fmt.Sprintf("docker logs -f %s", containerID)

	if isLocal(server) {
		return m.streamLocalLogs(ctx, logCmd, artifactUID, logc, errc)
	}

	return m.streamRemoteLogs(ctx, server, logCmd, artifactUID, logc, errc)
}

func (m *VmManager) streamLocalLogs(ctx context.Context, command, artifactUID string, logc chan<- *types.ArtifactLogLine, errc chan<- error) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = cmd.Stdout // Merge stderr into stdout

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(pipe)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if scanner.Scan() {
				logc <- &types.ArtifactLogLine{
					ArtifactUID: artifactUID,
					Log:         scanner.Text(),
				}
			} else {
				return scanner.Err()
			}
		}
	}
}

func (m *VmManager) streamRemoteLogs(ctx context.Context, server *types.Server, command, artifactUID string, logc chan<- *types.ArtifactLogLine, errc chan<- error) error {
	reader, cleanup, err := m.runSSHStream(ctx, server, command)
	if err != nil {
		return err
	}
	defer cleanup()

	scanner := bufio.NewScanner(reader)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if scanner.Scan() {
				logc <- &types.ArtifactLogLine{
					ArtifactUID: artifactUID,
					Log:         scanner.Text(),
				}
			} else {
				return scanner.Err()
			}
		}
	}
}
