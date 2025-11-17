package environment

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Get(ctx context.Context, projectID int64, environmentUID int64) (*types.Environment, error) {
	return c.FindByUID(ctx, projectID, environmentUID)
}
