package volume

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Update(ctx context.Context, volume *types.Volume) (*types.Volume, error) {
	if err := c.sanitizeCreateInput(volume.ToInput()); err != nil {
		return nil, err
	}
	return c.volumeStore.Update(ctx, volume)
}

func (c *Controller) Detach(ctx context.Context, volume *types.Volume) (*types.Volume, error) {
	volume.ApplicaitonID = nil
	return c.volumeStore.Update(ctx, volume)
}
