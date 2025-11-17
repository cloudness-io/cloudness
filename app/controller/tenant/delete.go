package tenant

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) SoftDelete(ctx context.Context, tenant *types.Tenant) error {
	now := time.Now().UTC().UnixMilli()

	return c.tx.WithTx(ctx, func(ctx context.Context) error {
		if err := c.tenantStore.SoftDelete(ctx, tenant, now); err != nil {
			return err
		}

		if err := c.tenantMembershipStore.DeleteAll(ctx, tenant.ID); err != nil {
			return err
		}

		return c.projectCtrl.SoftDeleteInTenant(ctx, tenant.ID, now)
	})
}
