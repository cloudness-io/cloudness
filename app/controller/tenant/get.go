package tenant

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) GetRestrctions(ctx context.Context, tenant *types.Tenant) *types.TenantRestrictions {
	return c.configSvc.GetTenantRestrictions(tenant)
}
