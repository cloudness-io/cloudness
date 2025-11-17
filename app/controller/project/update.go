package project

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Update(ctx context.Context, projectID int64, in *CreateProjectInput) (*types.Project, error) {
	if err := c.sanitizeCreateInput(in); err != nil {
		return nil, err
	}

	project, err := c.findByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	project.Name = in.Name
	project.Description = in.Description

	return c.projectStore.Update(ctx, project)
}
