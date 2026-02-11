package environment

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/environment"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/venvironment"

	"github.com/rs/zerolog/log"
)

func HandleNew() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		render.Page(ctx, w, venvironment.AddEnvironmentPage())
	}
}

func HandleAdd(envCtrl *environment.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)

		in := new(environment.CreateEnvironmentInput)
		err := json.NewDecoder(r.Body).Decode(in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding json")
			render.ToastErrorMsg(ctx, w, "Invalid Request Body")
			return
		}

		env, err := envCtrl.Create(ctx, session, tenant, project, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("error creating project")
			render.ToastError(ctx, w, err)
			return
		}

		ctx = request.WithEnvironment(ctx, env)

		render.Redirect(w, routes.ProjectEnvironment)
	}
}
