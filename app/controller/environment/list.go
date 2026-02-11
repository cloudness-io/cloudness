package environment

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

func (c *Controller) List(ctx context.Context, tenantID, projectID int64) ([]*types.Environment, error) {
	filter := &types.EnvironmentFilter{
		TenantID:  &tenantID,
		ProjectID: &projectID,
		Sort:      enum.EnvironmentAttrSequence,
		Order:     enum.OrderDesc,
	}
	return c.environmentStore.List(ctx, filter)
}
