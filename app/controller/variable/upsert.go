package variable

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) upsert(ctx context.Context, variables *types.Variable) error {
	return c.variableStore.Upsert(ctx, variables)
}

func (c *Controller) UpsertMany(ctx context.Context, variables []*types.Variable) error {
	return c.variableStore.UpsertMany(ctx, variables)
}
