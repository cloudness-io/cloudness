package application

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/deployment"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vapplication"

	"github.com/rs/zerolog/log"
)

func HandleListDeployments(appCtrl *application.Controller, deploymentCtrl *deployment.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		application, _ := request.ApplicationFrom(ctx)

		deployments, err := deploymentCtrl.List(ctx, application.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing deployments")
			render.ToastError(ctx, w, err)
		}

		render.Page(ctx, w, vapplication.DeploymentsList(application, deployments))
	}
}
