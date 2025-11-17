package application

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

func HandleRedeploy(appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		application, _ := request.ApplicationFrom(ctx)

		deployment, err := appCtrl.Redeploy(ctx, session.Principal.DisplayName, application)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("error redeploying application")
			render.ToastError(ctx, w, err)
			return
		}

		render.Redirect(w, routes.DeploymentCtx(request.WithDeployment(ctx, deployment)))
	}
}

func HandleDeploy(appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		application, _ := request.ApplicationFrom(ctx)

		deployment, err := appCtrl.Deploy(ctx, session.Principal.DisplayName, enum.TriggerActionManual, application)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("error redeploying application")
			render.ToastError(ctx, w, err)
			return
		}

		render.Redirect(w, routes.DeploymentCtx(request.WithDeployment(ctx, deployment)))
	}
}
