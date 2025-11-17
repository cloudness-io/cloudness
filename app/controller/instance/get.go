package instance

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Get(ctx context.Context) (*types.Instance, error) {
	return c.instanceStore.Get(ctx)
}
