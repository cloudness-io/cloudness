package environment

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/environment"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vproject"

	"github.com/rs/zerolog/log"
)

func HandleGet(envCtrl *environment.Controller, appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)

		apps, err := appCtrl.List(ctx, env.TenantID, env.ProjectID, env.ID)
		if err != nil {
			log.Error().Err(err).Msg("error listing applications")
			render.ToastError(ctx, w, err)
			return
		}

		render.Page(ctx, w, vproject.Overview(project, env, apps))
	}
}
