package server

import (
	"context"
	"errors"

	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/app/services/dns"
	"github.com/cloudness-io/cloudness/app/services/manager"
	"github.com/cloudness-io/cloudness/app/services/proxy"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/helpers"
	dbStore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

type Controller struct {
	tx            dbtx.Transactor
	configSvc     *config.Service
	dnsSvc        *dns.Service
	proxySvc      *proxy.Service
	serverStore   store.ServerStore
	instanceStore store.InstanceStore
	factory       manager.ManagerFactory
}

func NewController(
	tx dbtx.Transactor,
	configSvc *config.Service,
	dnsSvc *dns.Service,
	proxySvc *proxy.Service,
	serverStore store.ServerStore,
	instanceStore store.InstanceStore,
	factory manager.ManagerFactory,
) *Controller {
	return &Controller{
		tx:            tx,
		configSvc:     configSvc,
		dnsSvc:        dnsSvc,
		proxySvc:      proxySvc,
		serverStore:   serverStore,
		instanceStore: instanceStore,
		factory:       factory,
	}
}

func (c *Controller) resolveServerManager(server *types.Server) (manager.ServerManager, error) {
	return c.factory.GetServerManager(server)
}

func (c *Controller) Init(ctx context.Context) (*types.Server, error) {
	server, err := c.serverStore.Find(ctx, 1)
	if err != nil && !errors.Is(err, dbStore.ErrResourceNotFound) {
		return nil, err
	}

	if server != nil {
		return server, nil
	}

	defaults := c.configSvc.GetKubeServerDefaults()

	server = &types.Server{
		UID:                           helpers.GenerateUID(),
		Name:                          "Primary",
		Description:                   "Primary cluster",
		Type:                          enum.ServerTypeK8s,
		VolumeSupportsOnlineExpansion: defaults.DefaultVolumeSupportsOnlineExpansion,
		BuildEnabled:                  true,
		IsBuildServer:                 false,
		PollingInterval:               1,
		MaxConcurrentBuilds:           3,
		MaxCPUPerBuild:                1,
		MaxMemoryPerBuild:             2,
		VolumeMinSize:                 1,
	}

	manager, err := c.resolveServerManager(server)
	if err != nil {
		return nil, err
	}

	ip, err := manager.GetIP(ctx, server)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Unable to find ip address, skipping....")
	}

	server.IPV4 = ip

	return c.serverStore.Create(ctx, server)
}

func (c *Controller) findByID(ctx context.Context, serverID int64) (*types.Server, error) {
	return c.serverStore.Find(ctx, serverID)
}
