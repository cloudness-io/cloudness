package environment

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) SoftDelete(ctx context.Context, environment *types.Environment) error {
	now := time.Now().UTC().UnixMilli()

	return c.tx.WithTx(ctx, func(ctx context.Context) error {
		if err := c.environmentStore.SoftDelete(ctx, environment, now); err != nil {
			return err
		}

		if err := c.appCtrl.SoftDeleteInEnvironment(ctx, environment.ID, now); err != nil {
			return err
		}
		return nil
	})
}

func (c *Controller) SoftDeleteInProject(ctx context.Context, projectID, now int64) error {
	envs, err := c.environmentStore.List(ctx, &types.EnvironmentFilter{
		ProjectID: &projectID,
	})
	if err != nil {
		return err
	}

	for _, env := range envs {
		if err := c.environmentStore.SoftDelete(ctx, env, now); err != nil {
			return err
		}

		if err := c.appCtrl.SoftDeleteInEnvironment(ctx, env.ID, now); err != nil {
			return err
		}
	}

	return nil
}
