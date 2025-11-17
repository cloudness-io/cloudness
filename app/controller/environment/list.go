package environment

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) List(ctx context.Context, tenantID, projectID int64) ([]*types.Environment, error) {
	filter := &types.EnvironmentFilter{
		TenantID:  &tenantID,
		ProjectID: &projectID,
	}
	return c.environmentStore.List(ctx, filter)
}
