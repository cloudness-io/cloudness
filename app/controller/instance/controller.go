package instance

import (
	"context"
	"errors"

	"github.com/cloudness-io/cloudness/app/services/dns"
	"github.com/cloudness-io/cloudness/app/services/manager"
	"github.com/cloudness-io/cloudness/app/services/proxy"
	"github.com/cloudness-io/cloudness/app/store"
	dbStore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

type Controller struct {
	tx             dbtx.Transactor
	instanceStore  store.InstanceStore
	serverStore    store.ServerStore
	principalStore store.PrincipalStore
	dnsSvc         *dns.Service
	proxySvc       *proxy.Service
	factory        manager.ManagerFactory
}

func NewController(
	tx dbtx.Transactor,
	instanceStore store.InstanceStore,
	serverStore store.ServerStore,
	principalStore store.PrincipalStore,
	dnsSvc *dns.Service,
	proxySvc *proxy.Service,
	factory manager.ManagerFactory,
) *Controller {
	return &Controller{
		tx:             tx,
		instanceStore:  instanceStore,
		serverStore:    serverStore,
		principalStore: principalStore,
		dnsSvc:         dnsSvc,
		proxySvc:       proxySvc,
		factory:        factory,
	}
}

func (c *Controller) Init(ctx context.Context, server *types.Server) (*types.Instance, error) {
	instance, err := c.instanceStore.Get(ctx)
	if err != nil && !errors.Is(err, dbStore.ErrResourceNotFound) {
		return nil, err
	}
	if instance != nil {
		return instance, nil
	}

	instance = &types.Instance{
		PublicIPV4:            server.IPV4,
		PublicIPV6:            server.IPV6,
		UpdateEnabled:         true,
		DNSValidationEnabled:  true,
		DNSServers:            "8.8.8.8,1.1.1.1",
		UserSignupEnabled:     false,
		DemoUserEnabled:       false,
		RegistryEnabled:       true,
		RegistrySize:          1,
		RegistryMirrorEnabled: false,
		RegistryMirrorSize:    1,
	}

	instance, err = c.instanceStore.Create(ctx, instance)
	if err != nil {
		return nil, err
	}
	factory, err := c.factory.GetServerManager(server)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get server manager, skipping....")
		return instance, nil
	}

	if err := factory.AddOrUpdateRegistry(ctx, server, registryKey, instance.RegistrySize, true); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to add registry, skipping....")
	}

	return instance, nil
}

func (c *Controller) IsUserSignupAllowed(ctx context.Context) (bool, error) {
	usrCount, err := c.principalStore.CountUsers(ctx)
	if err != nil {
		return false, err
	}

	if usrCount == 0 {
		return true, nil
	}

	instance, err := c.Get(ctx)
	if err != nil {
		return false, err
	}

	return instance.UserSignupEnabled, nil
}
