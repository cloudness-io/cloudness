package server

import (
	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/app/services/dns"
	"github.com/cloudness-io/cloudness/app/services/manager"
	"github.com/cloudness-io/cloudness/app/services/proxy"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(
	tx dbtx.Transactor,
	configSvc *config.Service,
	dnsSvc *dns.Service,
	proxySvc *proxy.Service,
	serverStore store.ServerStore,
	instanceStore store.InstanceStore,
	factory manager.ManagerFactory,
) *Controller {
	return NewController(
		tx,
		configSvc,
		dnsSvc,
		proxySvc,
		serverStore,
		instanceStore,
		factory,
	)
}
