package favorite

import (
	"context"
	"errors"

	"github.com/cloudness-io/cloudness/store"
)

func (c *Controller) IsFavorite(ctx context.Context, userID, applicationID int64) (bool, error) {
	fav, err := c.favStore.Find(ctx, userID, applicationID)
	if err != nil {
		if errors.Is(err, store.ErrResourceNotFound) {
			return false, nil
		}
		return false, err
	}

	return fav != nil, err
}
