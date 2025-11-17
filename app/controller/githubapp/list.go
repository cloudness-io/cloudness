package githubapp

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) List(ctx context.Context, tenantID, projectID int64) ([]*types.GithubApp, error) {
	return c.ghAppSvc.List(ctx, tenantID, projectID)
}

func (c *Controller) ListConfigured(ctx context.Context, tenantID, projectID int64) ([]*types.GithubApp, error) {
	ghApps, err := c.List(ctx, tenantID, projectID)
	if err != nil {
		return nil, err
	}

	if len(ghApps) == 0 {
		return nil, nil
	}

	configured := []*types.GithubApp{}
	for _, ghApp := range ghApps {
		if ghApp.AppID != 0 {
			configured = append(configured, ghApp)
		}
	}
	return configured, nil
}
