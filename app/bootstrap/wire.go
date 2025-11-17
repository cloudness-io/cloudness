package bootstrap

import (
	"github.com/cloudness-io/cloudness/app/controller/auth"
	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/controller/server"
	"github.com/cloudness-io/cloudness/app/controller/template"
	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/types"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideBootstrap,
)

func ProvideBootstrap(
	config *types.Config,
	instanceCtrl *instance.Controller,
	serverCtrl *server.Controller,
	authCtrl *auth.Controller,
	userCtrl *user.Controller,
	templateCtrl *template.Controller,
) Bootstrap {
	return System(config, instanceCtrl, serverCtrl, authCtrl, userCtrl, templateCtrl)
}
