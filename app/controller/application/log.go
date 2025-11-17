package application

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) TailLog(ctx context.Context, app *types.Application) (<-chan *types.ArtifactLogLine, <-chan error, error) {
	server, err := c.serverCtrl.Get(ctx)
	if err != nil {
		return nil, nil, err
	}

	serverMgr, err := c.manager.GetServerManager(server)
	if err != nil {
		return nil, nil, err
	}
	return serverMgr.TailLogs(ctx, server, app)
}
