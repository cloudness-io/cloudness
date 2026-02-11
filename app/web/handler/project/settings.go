package project

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vproject"

	"github.com/rs/zerolog/log"
)

func HandleGetSettingsGeneral(projectCtrl *project.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		project, _ := request.ProjectFrom(ctx)

		render.Page(ctx, w, vproject.ProjectGeneral(project))
	}
}

func HandleUpdateSettingsGeneral(projectCtrl *project.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		p, _ := request.ProjectFrom(ctx)

		in := new(project.CreateProjectInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding json")
			render.ToastErrorMsg(ctx, w, "Invalid request body")
			return
		}

		p, err := projectCtrl.Update(ctx, p.ID, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating project")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		render.Redirect(w, routes.ProjectCtx(ctx)+"/"+routes.ProjectSettings)
	}
}
