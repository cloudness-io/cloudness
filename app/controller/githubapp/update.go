package githubapp

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) CompleteManifest(ctx context.Context, ghApp *types.GithubApp, code string) error {
	return c.ghAppSvc.CompleteManifest(ctx, ghApp, code)
}

func (c *Controller) CompleteInstallation(ctx context.Context, ghApp *types.GithubApp, installationID int64) error {
	return c.ghAppSvc.CompleteInstallation(ctx, ghApp, installationID)
}
