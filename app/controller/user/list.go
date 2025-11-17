package user

import "context"

func (c *Controller) CountTotal(ctx context.Context) (int64, error) {
	return c.principalStore.CountUsers(ctx)
}
