package cleanup

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func (j *deletedArtifactsJob) cleanupProjects(ctx context.Context, tenantID, deletedBeforeOrAt *int64, tracker *purgeTracker) error {
	for {
		filter := &types.ProjectFilter{
			TenantID: tenantID,
			ListQueryFilter: types.ListQueryFilter{
				Pagination: types.Pagination{
					Page: 1,
					Size: maxRetrivalItems,
				},
			},
			DeletedBeforeOrAt: deletedBeforeOrAt,
		}

		projectsToCleanup, err := j.projectStore.List(ctx, filter)
		if err != nil {
			return fmt.Errorf("failed to list projects to be cleaned up: %w", err)
		}

		if len(projectsToCleanup) == 0 {
			break
		}
		log.Ctx(ctx).Info().Msgf("found %d projects to be cleaned up", len(projectsToCleanup))

		for _, project := range projectsToCleanup {
			if err := j.projectStore.Purge(ctx, project.ID, project.Deleted); err != nil {
				return fmt.Errorf("failed to purge project %d: %w", project.ID, err)
			}

			tracker.project++
		}
	}

	return nil
}
