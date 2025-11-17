package template

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) List(ctx context.Context) ([]*types.Template, error) {
	return c.templateStore.List(ctx)
}

func (c *Controller) ListDatabase(ctx context.Context) ([]*types.Template, error) {
	return c.templateStore.ListBySlugs(ctx, []string{"postgres", "mysql", "valkey", "redis"})
}
