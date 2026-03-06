package vm

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func (m *VmManager) DeleteResources(ctx context.Context, server *types.Server, namespace, identifier string) error {
	// Remove all containers matching the app identifier on the given network (namespace).
	filter := fmt.Sprintf("--filter label=%s --filter label=%s", managedByFilter(), appFilter(identifier))
	cmd := fmt.Sprintf("docker ps -aq %s | xargs -r docker rm -f", filter)

	if _, err := m.runCmd(ctx, server, cmd); err != nil {
		log.Ctx(ctx).Warn().Err(err).Str("identifier", identifier).Msg("vm: delete resources")
	}

	return nil
}

func (m *VmManager) DeleteNamespace(ctx context.Context, server *types.Server, namespace string) error {
	// A namespace in the VM world is a Docker network.
	cmd := fmt.Sprintf("docker network rm %s 2>/dev/null || true", namespace)
	_, err := m.runCmd(ctx, server, cmd)
	return err
}

func (m *VmManager) DeleteVolume(ctx context.Context, server *types.Server, volume *types.Volume) error {
	cmd := fmt.Sprintf("docker volume rm %s 2>/dev/null || true", volume.GetIdentifierStr())
	_, err := m.runCmd(ctx, server, cmd)
	return err
}

func (m *VmManager) DeleteApplication(ctx context.Context, server *types.Server, app *types.Application) error {
	identifier := app.GetIdentifierStr()
	filter := fmt.Sprintf("--filter label=%s --filter label=%s", managedByFilter(), appFilter(identifier))

	// 1. Stop and remove all containers for this application.
	cmd := fmt.Sprintf("docker ps -aq %s | xargs -r docker rm -f", filter)
	if _, err := m.runCmd(ctx, server, cmd); err != nil {
		log.Ctx(ctx).Warn().Err(err).Str("app", identifier).Msg("vm: remove containers")
	}

	// 2. Remove associated Docker volumes (named with the app identifier prefix).
	cmd = fmt.Sprintf("docker volume ls -q --filter label=%s | xargs -r docker volume rm", appFilter(identifier))
	if _, err := m.runCmd(ctx, server, cmd); err != nil {
		log.Ctx(ctx).Warn().Err(err).Str("app", identifier).Msg("vm: remove volumes")
	}

	// 3. Remove Traefik dynamic config files for this application.
	cmd = fmt.Sprintf("rm -f %s/%s-*.yml", traefikDynamicConfigDir, identifier)
	if _, err := m.runCmd(ctx, server, cmd); err != nil {
		log.Ctx(ctx).Warn().Err(err).Str("app", identifier).Msg("vm: remove traefik configs")
	}

	return nil
}
