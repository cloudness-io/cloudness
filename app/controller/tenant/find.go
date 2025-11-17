package tenant

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Find(ctx context.Context, tenantID int64) (*types.Tenant, error) {
	tenant, err := c.findByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

func (c *Controller) FindByUID(ctx context.Context, tenantUID int64) (*types.Tenant, error) {
	return c.findByUID(ctx, tenantUID)
}
