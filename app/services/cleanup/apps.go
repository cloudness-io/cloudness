package cleanup

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func (j *deletedArtifactsJob) cleanupApps(ctx context.Context, tenantID, projectID, envID, deletedBeforeOrAt *int64, tracker *purgeTracker) error {
	for {
		filter := &types.ApplicationFilter{
			TenantID:      tenantID,
			ProjectID:     projectID,
			EnvironmentID: envID,
			ListQueryFilter: types.ListQueryFilter{
				Pagination: types.Pagination{
					Page: 1,
					Size: maxRetrivalItems,
				},
			},
			DeletedBeforeOrAt: deletedBeforeOrAt,
		}

		appsToCleanup, err := j.appStore.List(ctx, filter)
		if err != nil {
			return fmt.Errorf("failed to list apps to be cleaned up: %w", err)
		}

		if len(appsToCleanup) == 0 {
			break
		}
		log.Ctx(ctx).Info().Msgf("found %d apps to be cleaned up", len(appsToCleanup))

		for _, app := range appsToCleanup {
			if err := j.purgeApplicationArtifacts(ctx, app); err != nil {
				return fmt.Errorf("failed to purge artifacts artifact %d: %w", app.ID, err)
			}
			if err := j.appStore.Purge(ctx, app.ID, app.Deleted); err != nil {
				return fmt.Errorf("failed to purge application %d: %w", app.ID, err)
			}

			tracker.app++
		}
	}

	return nil
}

func (j *deletedArtifactsJob) purgeApplicationArtifacts(ctx context.Context, app *types.Application) error {
	server, manager, err := j.getServerManager(ctx, app.ServerID)
	if err != nil {
		return err
	}

	return manager.DeleteApplication(ctx, server, app)
}
