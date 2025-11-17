package environment

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/environment"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vproject"

	"github.com/rs/zerolog/log"
)

func HandleUpdate(envCtrl *environment.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		env, _ := request.EnvironmentFrom(ctx)

		in := new(environment.CreateEnvironmentInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding json")
			render.ToastErrorMsg(ctx, w, "Invalid request body")
			return
		}

		env, err := envCtrl.Update(ctx, env.ID, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating environment")
			render.ToastError(ctx, w, err)
			return
		}
		ctx = request.WithEnvironment(ctx, env)

		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		envs, err := envCtrl.List(ctx, tenant.ID, project.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing environments")
			render.Redirect(w, routes.ProjectCtx(ctx)+"/"+routes.ProjectEnvironment)
			return
		}

		render.Page(ctx, w, vproject.ListEnvs(project, envs, env))
		render.ToastSuccess(ctx, w, "Environment updated successfully")
	}
}
