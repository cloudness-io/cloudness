package application

import (
	"context"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/environment"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vapplication"
	"github.com/cloudness-io/cloudness/app/web/views/components/vproject"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func HandleList(envCtrl *environment.Controller, appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RenderAppList(r.Context(), w, r, nil, envCtrl, appCtrl)
	}
}

func RenderAppList(ctx context.Context, w http.ResponseWriter, r *http.Request, envs []*types.Environment, envCtrl *environment.Controller, appCtrl *application.Controller) {
	tenant, _ := request.TenantFrom(ctx)
	project, _ := request.ProjectFrom(ctx)
	env, _ := request.EnvironmentFrom(ctx)
	target := request.TargetElementFrom(ctx)

	apps, err := appCtrl.List(ctx, tenant.ID, project.ID, env.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Error listing apps")
		render.ToastError(ctx, w, err)
		return
	}

	if request.HxIndicatorFrom(ctx) && target != routes.TargetMain {
		render.HTMLWithBreadCrumb(ctx, w, vapplication.List(env, apps))
		return
	}

	if len(envs) == 0 {
		//is not hx request, should load full page
		envs, err = envCtrl.List(ctx, tenant.ID, project.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing environments")
			render.ToastError(ctx, w, err)
			return
		}
	}
	render.Page(ctx, w, vproject.Overview(project, envs, env, apps))
}
