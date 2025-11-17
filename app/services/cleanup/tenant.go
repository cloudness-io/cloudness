package cleanup

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func (j *deletedArtifactsJob) cleanupTenant(ctx context.Context, deletedBeforeOrAt *int64, tracker *purgeTracker) error {
	for {
		filter := &types.TenantFilter{
			ListQueryFilter: types.ListQueryFilter{
				Pagination: types.Pagination{
					Page: 1,
					Size: maxRetrivalItems,
				},
			},
			DeletedBeforeOrAt: deletedBeforeOrAt,
		}

		tenantsToCleanup, err := j.tenantStore.List(ctx, filter)
		if err != nil {
			return fmt.Errorf("failed to list tenants to be cleaned up: %w", err)
		}

		if len(tenantsToCleanup) == 0 {
			break
		}
		log.Ctx(ctx).Info().Msgf("found %d tenants to be cleaned up", len(tenantsToCleanup))

		for _, tenant := range tenantsToCleanup {
			if err := j.cleanupProjects(ctx, &tenant.ID, nil, tracker); err != nil {
				return err
			}

			if err := j.tenantStore.Purge(ctx, tenant.ID, tenant.Deleted); err != nil {
				return fmt.Errorf("failed to purge tenant %d: %w", tenant.ID, err)
			}

			tracker.tenant++
		}
	}

	return nil
}
