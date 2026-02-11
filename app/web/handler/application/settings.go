package application

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/githubapp"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vapplication"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func HandleGetSettings(appCtrl *application.Controller, ghAppCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if err := renderSettings(ctx, w, appCtrl, ghAppCtrl); err != nil {
			return
		}
	}
}

func HandleUpdateSettings(appCtrl *application.Controller, ghAppCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		app, _ := request.ApplicationFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		environment, _ := request.EnvironmentFrom(ctx)
		session, _ := request.AuthSessionFrom(ctx)

		in := new(types.ApplicationInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Error().Err(err).Msg("Invalid request body")
			render.ToastErrorMsg(ctx, w, "Invalid request body")
			return
		}

		app, err := appCtrl.Update(ctx, session, tenant, project, environment, app, in)
		if err != nil {
			log.Error().Err(err).Msg("Error updating application")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		ctx = request.WithApplication(ctx, app)
		if err := renderSettings(ctx, w, appCtrl, ghAppCtrl); err != nil {
			return
		}
		render.ToastSuccess(ctx, w, "Applicaiton updated successfully")
	}
}

func renderSettings(ctx context.Context, w http.ResponseWriter, appCtrl *application.Controller, ghAppCtrl *githubapp.Controller) error {
	app, _ := request.ApplicationFrom(ctx)
	tenant, _ := request.TenantFrom(ctx)
	project, _ := request.ProjectFrom(ctx)

	restrctions := appCtrl.GetRestrictions(ctx, tenant)
	var ghApp *types.GithubApp
	var err error
	if app.GetGithubAppID() > 0 {
		ghApp, err = ghAppCtrl.Find(ctx, app.TenantID, app.ProjectID, app.GetGithubAppID())
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting github app")
			render.ToastError(ctx, w, err)
			return err
		}
	}

	render.Page(ctx, w, vapplication.Settings(project, app, ghApp, restrctions))
	return nil
}
