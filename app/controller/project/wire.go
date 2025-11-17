package project

import (
	"github.com/cloudness-io/cloudness/app/controller/environment"
	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/app/sse"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProviderController,
)

func ProviderController(
	tx dbtx.Transactor,
	configSvc *config.Service,
	userCtrl *user.Controller,
	envCtrl *environment.Controller,
	projectStore store.ProjectStore,
	projectMembershipStore store.ProjectMembershipStore,
	tenantMembershipStore store.TenantMembershipStore,
	sseStremer sse.Streamer,
) *Controller {
	return NewController(tx,
		configSvc, userCtrl, envCtrl,
		projectStore, projectMembershipStore,
		tenantMembershipStore,
		sseStremer,
	)
}
