package cleanup

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func (j *deletedArtifactsJob) cleanupEnviornments(ctx context.Context, tenantID, projectID, deletedBeforeOrAt *int64, tracker *purgeTracker) error {
	for {
		filter := &types.EnvironmentFilter{
			TenantID:  tenantID,
			ProjectID: projectID,
			ListQueryFilter: types.ListQueryFilter{
				Pagination: types.Pagination{
					Page: 1,
					Size: maxRetrivalItems,
				},
			},
			DeletedBeforeOrAt: deletedBeforeOrAt,
		}

		envsToCleanup, err := j.envStore.List(ctx, filter)
		if err != nil {
			return fmt.Errorf("failed to list environments to be cleaned up: %w", err)
		}

		if len(envsToCleanup) == 0 {
			break
		}
		log.Ctx(ctx).Info().Msgf("found %d environments to be cleaned up", len(envsToCleanup))

		for _, env := range envsToCleanup {
			if err := j.purgeEnvironmentArtifacts(ctx, env); err != nil {
				return err
			}

			if err := j.envStore.Purge(ctx, env.ID, env.Deleted); err != nil {
				return fmt.Errorf("failed to purge environment %d: %w", env.ID, err)
			}

			tracker.env++
		}
	}

	return nil
}

func (j *deletedArtifactsJob) purgeEnvironmentArtifacts(ctx context.Context, env *types.Environment) error {
	servers, err := j.serverStore.List(ctx)
	if err != nil {
		return err
	}
	for _, server := range servers {
		manager, err := j.getServerManagerForServer(ctx, server)
		if err != nil {
			return err
		}

		if err := manager.DeleteNamespace(ctx, server, env.Slug); err != nil {
			return err
		}
	}
	return nil
}
