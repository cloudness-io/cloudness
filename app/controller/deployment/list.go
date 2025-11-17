package deployment

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) List(ctx context.Context, applicationID int64) ([]*types.Deployment, error) {
	return c.deploymentStore.List(ctx, applicationID)
}
