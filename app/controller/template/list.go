package template

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

var databaseSlugs = []string{"postgres", "mysql", "valkey", "redis"}

func (c *Controller) List(ctx context.Context) ([]*types.Template, error) {
	return c.templateStore.List(ctx)
}

func (c *Controller) ListDatabase(ctx context.Context) ([]*types.Template, error) {
	return c.templateStore.ListBySlugs(ctx, databaseSlugs)
}

func (c *Controller) ListTemplates(ctx context.Context) ([]*types.Template, error) {
	return c.templateStore.ListNotInSlugs(ctx, databaseSlugs)
}
