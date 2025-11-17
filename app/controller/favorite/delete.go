package favorite

import "context"

func (c *Controller) Delete(ctx context.Context, userID, applicationID int64) error {
	return c.favStore.Delete(ctx, userID, applicationID)
}
