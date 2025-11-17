package environment

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Update(ctx context.Context, environmentID int64, in *CreateEnvironmentInput) (*types.Environment, error) {
	if err := c.sanitizeCreateInput(in); err != nil {
		return nil, err
	}

	env, err := c.findEnvironmentByID(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	env.Name = in.Name

	return c.environmentStore.Update(ctx, env)
}
