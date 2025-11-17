package deployment

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Get(ctx context.Context, applicationID, deploymentUID int64) (*types.Deployment, error) {
	return c.findByUID(ctx, applicationID, deploymentUID)
}
