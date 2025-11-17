package project

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog/log"
)

func HandleAdd(projectCtrl *project.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)

		in := new(project.CreateProjectInput)
		err := json.NewDecoder(r.Body).Decode(in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding json")
			render.ToastErrorMsg(ctx, w, "Invalid Request Body")
			return
		}

		project, err := projectCtrl.Create(ctx, session, tenant, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("error creating project")
			render.ToastError(ctx, w, err)
			return
		}
		ctx = request.WithProject(ctx, project)

		render.Redirect(w, routes.ProjectCtx(ctx))
	}
}
