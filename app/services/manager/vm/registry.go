package vm

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/types"
)

func (m *VmManager) AddOrUpdateRegistry(ctx context.Context, server *types.Server, name string, size int64, enableNodePort bool) error {
	portMapping := fmt.Sprintf("-p 127.0.0.1:%d:%d", defaultRegistryPort, defaultRegistryPort)
	if enableNodePort {
		portMapping = fmt.Sprintf("-p %d:%d", defaultRegistryPort, defaultRegistryPort)
	}

	// Create a Docker volume for registry data.
	volCmd := fmt.Sprintf("docker volume create %s-data", name)
	if _, err := m.runCmd(ctx, server, volCmd); err != nil {
		return fmt.Errorf("failed to create registry volume: %w", err)
	}

	// Stop existing registry container if any.
	stopCmd := fmt.Sprintf("docker rm -f %s 2>/dev/null || true", name)
	if _, err := m.runCmd(ctx, server, stopCmd); err != nil {
		return fmt.Errorf("failed to stop existing registry: %w", err)
	}

	// Run the registry container.
	cmd := fmt.Sprintf(
		"docker run -d --name %s --restart unless-stopped "+
			"%s "+
			"-v %s-data:/var/lib/registry "+
			"--label %s=%s "+
			"%s",
		name,
		portMapping,
		name,
		labelManagedBy, labelManagedByValue,
		defaultRegistryImage,
	)

	_, err := m.runCmd(ctx, server, cmd)
	if err != nil {
		return fmt.Errorf("failed to start registry: %w", err)
	}

	return nil
}

func (m *VmManager) RemoveRegistry(ctx context.Context, server *types.Server, name string) error {
	// Stop and remove the registry container.
	cmd := fmt.Sprintf("docker rm -f %s 2>/dev/null || true", name)
	if _, err := m.runCmd(ctx, server, cmd); err != nil {
		return fmt.Errorf("failed to remove registry container: %w", err)
	}

	// Remove the data volume.
	cmd = fmt.Sprintf("docker volume rm %s-data 2>/dev/null || true", name)
	_, err := m.runCmd(ctx, server, cmd)
	return err
}
