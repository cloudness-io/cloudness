package project

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/environment"
	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vproject"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func HandleGet(projectCtrl *project.Controller, envCtrl *environment.Controller, appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		var env *types.Environment
		var apps []*types.Application

		envs, err := envCtrl.List(ctx, tenant.ID, project.ID)
		if err != nil {
			log.Error().Err(err).Msg("error listing environments")
			render.ToastError(ctx, w, err)
			return
		}
		if len(envs) > 0 {
			env = envs[0]
			// ctx = request.WithTargetElement(ctx, routes.TargetMain)
			// w.Header().Set("HX-Push-Url", routes.EnvironmentCtx(ctx)+routes.EnvironmentApplication)
			// handlerapplication.RenderAppList(ctx, w, r, envs, envCtrl, appCtrl)
			// return
			render.RedirectWithRefresh(w, routes.EnvironmentCtxUID(ctx, env.UID))
			return
		}

		if env != nil {
			apps, err = appCtrl.List(ctx, tenant.ID, project.ID, env.ID)
			if err != nil {
				log.Error().Err(err).Msg("error listing applications")
				render.ToastError(ctx, w, err)
				return
			}
		}

		render.Page(ctx, w, vproject.Overview(project, env, apps))
	}
}
