package auth

import (
	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(
	tx dbtx.Transactor,
	instanceCtrl *instance.Controller,
	userCtrl *user.Controller,
	tenantCtrl *tenant.Controller,
	tokenStore store.TokenStore,
	authSettingStore store.AuthSettingsStore,
) *Controller {
	return New(
		tx,
		instanceCtrl,
		userCtrl,
		tenantCtrl,
		tokenStore,
		authSettingStore,
	)
}
