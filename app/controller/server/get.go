package server

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

// TODO: end goal is to remove this when we start supporting multiple servers
func (c *Controller) Get(ctx context.Context) (*types.Server, error) {
	return c.findByID(ctx, 1)
}

func (c *Controller) FindByID(ctx context.Context, serverID int64) (*types.Server, error) {
	return c.findByID(ctx, serverID)
}
