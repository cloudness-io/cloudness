package favorite

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) List(ctx context.Context, userID, tenantID int64) ([]*types.FavoriteDTO, error) {
	return c.favStore.List(ctx, userID, tenantID)
}
