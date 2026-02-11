package connections

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/githubapp"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vsource/vgithubapp"

	"github.com/rs/zerolog/log"
)

func HandleNewGithubApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		render.Page(ctx, w, vgithubapp.AddGithubAppPage())
	}
}

func HandleAddGithubApp(ghCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		var projectID int64 = 0

		in := new(githubapp.CreateGithubAppInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		if project, ok := request.ProjectFrom(ctx); ok {
			projectID = project.ID
		}

		ghApp, err := ghCtrl.Create(ctx, session, tenant.ID, int64(projectID), in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error creating github app")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		render.RedirectWithRefresh(w, routes.ProjectConnectionGithubUIDCtx(ctx, ghApp.UID))
	}
}
