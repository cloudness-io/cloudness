package template

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) FindByID(ctx context.Context, id int64) (*types.Template, error) {
	return c.templateStore.Find(ctx, id)
}
