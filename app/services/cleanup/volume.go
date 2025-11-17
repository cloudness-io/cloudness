package cleanup

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func (j *deletedArtifactsJob) cleanupVolumes(ctx context.Context, tenantID, projectID, envID, deletedBeforeOrAt *int64, tracker *purgeTracker) error {
	for {
		filter := &types.VolumeFilter{
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

		volumesToCleanup, err := j.volumeStore.List(ctx, filter)
		if err != nil {
			return fmt.Errorf("failed to list volumes to be cleaned up: %w", err)
		}

		if len(volumesToCleanup) == 0 {
			break
		}
		log.Ctx(ctx).Info().Msgf("found %d volumes to be cleaned up", len(volumesToCleanup))

		for _, volume := range volumesToCleanup {
			if err := j.purgeVolumeArtifact(ctx, volume); err != nil {
				return fmt.Errorf("failed to purge volume artifact %d: %w", volume.ID, err)
			}
			if err := j.volumeStore.Purge(ctx, volume.ID, volume.Deleted); err != nil {
				return fmt.Errorf("failed to purge volume %d: %w", volume.ID, err)
			}

			tracker.volume++
		}
	}

	return nil
}

func (j *deletedArtifactsJob) purgeVolumeArtifact(ctx context.Context, volume *types.Volume) error {
	server, manager, err := j.getServerManager(ctx, volume.ServerID)
	if err != nil {
		return err
	}

	return manager.DeleteVolume(ctx, server, volume)
}
