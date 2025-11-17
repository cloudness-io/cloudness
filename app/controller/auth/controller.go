package auth

import (
	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
)

type Controller struct {
	tx               dbtx.Transactor
	instanceCtrl     *instance.Controller
	userCtrl         *user.Controller
	tenantCtrl       *tenant.Controller
	tokenStore       store.TokenStore
	authSettingStore store.AuthSettingsStore
}

func New(
	tx dbtx.Transactor,
	instanceCtrl *instance.Controller,
	userCtrl *user.Controller,
	tenantCtrl *tenant.Controller,
	tokenStore store.TokenStore,
	authSettingStore store.AuthSettingsStore,
) *Controller {
	return &Controller{
		tx:               tx,
		instanceCtrl:     instanceCtrl,
		userCtrl:         userCtrl,
		tenantCtrl:       tenantCtrl,
		tokenStore:       tokenStore,
		authSettingStore: authSettingStore,
	}
}
