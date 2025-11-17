package application

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Get(ctx context.Context, tenantID, projectID, environmentID, applicationUID int64) (*types.Application, error) {
	return c.findByUID(ctx, tenantID, projectID, environmentID, applicationUID)
}
