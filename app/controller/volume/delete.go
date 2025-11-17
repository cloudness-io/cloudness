package volume

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) SoftDelete(ctx context.Context, volume *types.Volume) error {
	now := time.Now().UTC().UnixMilli()
	return c.volumeStore.SoftDelete(ctx, volume, now)
}

func (c *Controller) SoftDeleteInApplication(ctx context.Context, appID, now int64) error {
	volumes, err := c.volumeStore.List(ctx, &types.VolumeFilter{
		ApplicationID: &appID,
	})
	if err != nil {
		return err
	}

	for _, volume := range volumes {
		if err := c.volumeStore.SoftDelete(ctx, volume, now); err != nil {
			return err
		}
	}

	return nil
}
