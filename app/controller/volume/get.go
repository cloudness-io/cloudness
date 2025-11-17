package volume

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Get(ctx context.Context, tenantID, projectID, envID, volumeUID int64) (*types.Volume, error) {
	return c.findByUID(ctx, tenantID, projectID, envID, volumeUID)
}
