package tenant

import (
	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProviderController,
)

func ProviderController(
	tx dbtx.Transactor,
	configSvc *config.Service,
	tenantStore store.TenantStore,
	tenantMembershipStore store.TenantMembershipStore,
	userCtrl *user.Controller,
	projectCtrl *project.Controller,
) *Controller {
	return NewController(tx, configSvc, tenantStore, tenantMembershipStore, userCtrl, projectCtrl)
}
