package favorite

import "context"

func (c *Controller) Add(ctx context.Context, userID, applicationID int64) error {
	return c.favStore.Add(ctx, userID, applicationID)
}
