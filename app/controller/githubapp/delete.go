package githubapp

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Delete(ctx context.Context, ghApp *types.GithubApp) error {
	return c.ghAppSvc.Delete(ctx, ghApp)
}
