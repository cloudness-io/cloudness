package application

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) GetLogs(ctx context.Context, app *types.Application) ([]*types.Artifact, error) {
	server, err := c.serverCtrl.Get(ctx)
	if err != nil {
		return nil, err
	}

	mgr, err := c.manager.GetServerManager(server)
	if err != nil {
		return nil, err
	}

	return mgr.ListArtifacts(ctx, server, app)
}
