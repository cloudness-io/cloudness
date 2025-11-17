package application

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) GetRestrictions(ctx context.Context, tenant *types.Tenant) *types.ApplicationRestrction {
	return c.configSvc.GetAppRestrictions(tenant)
}
