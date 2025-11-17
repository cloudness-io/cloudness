package application

import (
	"context"
	"errors"

	dbStore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) List(ctx context.Context, tenantID, projectID, environmentID int64) ([]*types.Application, error) {
	apps, err := c.applicationStore.List(ctx, &types.ApplicationFilter{
		TenantID:      &tenantID,
		ProjectID:     &projectID,
		EnvironmentID: &environmentID,
	})
	if err != nil && !errors.Is(err, dbStore.ErrResourceNotFound) {
		return nil, err
	}

	return apps, nil
}
