package router

import (
	"context"

	"github.com/cloudness-io/cloudness/app/auth/authn"
	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/auth"
	"github.com/cloudness-io/cloudness/app/controller/deployment"
	"github.com/cloudness-io/cloudness/app/controller/environment"
	"github.com/cloudness-io/cloudness/app/controller/favorite"
	"github.com/cloudness-io/cloudness/app/controller/githubapp"
	"github.com/cloudness-io/cloudness/app/controller/gitpublic"
	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/controller/logs"
	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/controller/server"
	"github.com/cloudness-io/cloudness/app/controller/template"
	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/app/controller/variable"
	"github.com/cloudness-io/cloudness/app/controller/volume"
	"github.com/cloudness-io/cloudness/types"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideRouter,
	// ProvideAPIHandler,
	ProvideWebHandler,
)

func ProvideRouter(
	// api APIHandler,
	web WebHandler,
) *Router {
	// return NewRouter(api, web)
	return NewRouter(web)
}

// func ProvideAPIHandler(
// 	appCtx context.Context,
// 	config *types.Config,
// 	authenticator authn.Authenticator,
// 	userCtrl *user.Controller,
// 	tenatCtrl *tenant.Controller,
// 	projectCtrl *project.Controller,
// 	appCtrl *app.Controller,
// 	spCtrl *sourceprovider.Controller,
// ) APIHandler {
// 	return NewAPIHandler(appCtx, config,
// 		authenticator, userCtrl,
// 		tenatCtrl, projectCtrl,
// 		appCtrl, spCtrl)
// }

func ProvideWebHandler(
	appCtx context.Context,
	config *types.Config,
	authenticator authn.Authenticator,
	instanceCtrl *instance.Controller,
	serverCtrl *server.Controller,
	userCtrl *user.Controller,
	tenatCtrl *tenant.Controller,
	projectCtrl *project.Controller,
	environmentCtrl *environment.Controller,
	authCtrl *auth.Controller,
	ghAppCtrl *githubapp.Controller,
	gitPublicCtrl *gitpublic.Controller,
	appCtrl *application.Controller,
	varCtrl *variable.Controller,
	deploymentCtrl *deployment.Controller,
	logsCtrl *logs.Controller,
	volumeCtrl *volume.Controller,
	templCtrl *template.Controller,
	favCtrl *favorite.Controller,
) WebHandler {
	return NewWebHandler(appCtx, config,
		authenticator,
		instanceCtrl, serverCtrl, userCtrl,
		tenatCtrl, projectCtrl,
		environmentCtrl, authCtrl,
		ghAppCtrl, gitPublicCtrl,
		appCtrl, varCtrl, deploymentCtrl,
		logsCtrl, volumeCtrl,
		templCtrl, favCtrl,
	)
}
