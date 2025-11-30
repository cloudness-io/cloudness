package instance

import (
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
	instanceStore store.InstanceStore,
	serverStore store.ServerStore,
	principalStore store.PrincipalStore,
	dnsSvc *dns.Service,
	proxySvc *proxy.Service,
	factory manager.ManagerFactory,
) *Controller {
	return NewController(tx, instanceStore, serverStore, principalStore, dnsSvc, proxySvc, factory)
}
