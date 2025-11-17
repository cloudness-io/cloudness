package volume

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) ListForApp(ctx context.Context, app *types.Application) ([]*types.Volume, error) {
	filter := &types.VolumeFilter{
		TenantID:      &app.TenantID,
		ProjectID:     &app.ProjectID,
		EnvironmentID: &app.EnvironmentID,
		ApplicationID: &app.ID,
	}

	return c.volumeStore.List(ctx, filter)
}

func (c *Controller) ListUnattached(ctx context.Context, tenantID, projectID, envID int64, serverID *int64) ([]*types.Volume, error) {
	filter := &types.VolumeFilter{
		TenantID:      &tenantID,
		ProjectID:     &projectID,
		EnvironmentID: &envID,
		ApplicationID: nil,
		ServerID:      serverID,
	}

	return c.volumeStore.List(ctx, filter)
}
