package cleanup

import (
	"context"

	"github.com/cloudness-io/cloudness/app/services/manager"
	"github.com/cloudness-io/cloudness/types"
)

func (j *deletedArtifactsJob) getServerManager(ctx context.Context, serverID int64) (*types.Server, manager.ServerManager, error) {
	server, err := j.serverStore.Find(ctx, serverID)
	if err != nil {
		return nil, nil, err
	}

	manager, err := j.getServerManagerForServer(ctx, server)
	if err != nil {
		return nil, nil, err
	}

	return server, manager, nil
}

func (j *deletedArtifactsJob) getServerManagerForServer(ctx context.Context, server *types.Server) (manager.ServerManager, error) {
	return j.serverFactory.GetServerManager(server)
}
