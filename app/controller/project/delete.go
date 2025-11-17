package project

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) SoftDelete(ctx context.Context, project *types.Project) error {
	now := time.Now().UTC().UnixMilli()

	return c.tx.WithTx(ctx, func(ctx context.Context) error {
		err := c.projectStore.SoftDelete(ctx, project, now)
		if err != nil {
			return nil
		}

		return c.envCtrl.SoftDeleteInProject(ctx, project.ID, now)
	})
}

func (c *Controller) SoftDeleteInTenant(ctx context.Context, tenantID, now int64) error {
	projects, err := c.projectStore.List(ctx, &types.ProjectFilter{
		TenantID: &tenantID,
	})
	if err != nil {
		return err
	}

	for _, project := range projects {
		if err := c.projectStore.SoftDelete(ctx, project, now); err != nil {
			return err
		}

		if err := c.envCtrl.SoftDeleteInProject(ctx, project.ID, now); err != nil {
			return err
		}
	}
	return nil
}
