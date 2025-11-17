package variable

import (
	"context"
	"errors"

	"github.com/cloudness-io/cloudness/store"
)

func (c *Controller) Delete(ctx context.Context, appID, varUID int64) error {
	return c.variableStore.Delete(ctx, appID, varUID)
}

func (c *Controller) DeleteByKey(ctx context.Context, appID int64, key string) error {
	err := c.variableStore.DeleteByKey(ctx, appID, key)
	if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
		return err
	}
	return nil
}
