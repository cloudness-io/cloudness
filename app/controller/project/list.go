package project

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) List(ctx context.Context, tenantID int64, principalID int64) ([]*types.Project, error) {
	return c.projectStore.List(ctx, &types.ProjectFilter{
		TenantID:    &tenantID,
		PrincipalID: &principalID,
	})
}
