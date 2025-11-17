package variable

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) ListInEnvironment(ctx context.Context, envID int64) ([]*types.Variable, error) {
	return c.variableStore.ListInEnvironment(ctx, envID)
}

func (c *Controller) ListDTO(ctx context.Context, envID, appID int64) (*types.VariableDTO, error) {
	allVars, err := c.ListInEnvironment(ctx, envID)
	if err != nil {
		return nil, err
	}
	return c.ToDTO(allVars, appID), nil
}

func (c *Controller) ListWithUpdate(ctx context.Context, envID, appID int64) (map[string]*types.Variable, error) {
	allVars, err := c.ListInEnvironment(ctx, envID)
	if err != nil {
		return nil, err
	}

	appVars, updatedVars := c.toListUpdate(allVars, appID)
	if len(updatedVars) > 0 {
		if err := c.UpsertMany(ctx, updatedVars); err != nil {
			return nil, err
		}
	}

	return appVars, nil
}
