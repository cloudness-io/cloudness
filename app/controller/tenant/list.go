package tenant

import (
	"context"
	"errors"

	"github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) ListMembership(ctx context.Context, userID int64) ([]*types.TenantMembershipUser, error) {
	tenants, err := c.tenantMembershipStore.List(ctx, userID)
	if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
		return nil, err
	}
	return tenants, nil
}

func (c *Controller) ListAllMembers(ctx context.Context, tenantID int64) ([]*types.TenantMembershipUser, error) {
	return c.tenantMembershipStore.ListByTenant(ctx, tenantID)
}
