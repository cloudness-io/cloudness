package githubapp

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Find(ctx context.Context, tenantID, projectID, githubAppID int64) (*types.GithubApp, error) {
	return c.ghAppSvc.Find(ctx, tenantID, projectID, githubAppID)
}

func (c *Controller) FindByUID(ctx context.Context, teanntID, projectID, githubAppUID int64) (*types.GithubApp, error) {
	return c.ghAppSvc.FindByUID(ctx, teanntID, projectID, githubAppUID)
}
