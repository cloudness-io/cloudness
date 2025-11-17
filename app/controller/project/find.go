package project

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Find(ctx context.Context, projectID int64) (*types.Project, error) {
	return c.findByID(ctx, projectID)
}

func (c *Controller) FindByUID(ctx context.Context, tenantID int64, projectUID int64) (*types.Project, error) {
	return c.findByUID(ctx, tenantID, projectUID)
}
