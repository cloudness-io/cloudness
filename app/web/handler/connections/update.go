package connections

import (
	"net/http"
	"strconv"

	"github.com/cloudness-io/cloudness/app/controller/githubapp"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog/log"
)

func HandleGithubCallback(ghCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		code := r.URL.Query().Get("code")
		// _ := r.URL.Query().Get("state")

		ghApp, _ := request.GithubAppFrom(ctx)

		err := ghCtrl.CompleteManifest(ctx, ghApp, code)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error completing github manifest")
			render.Error500(w, r)
			return
		}

		render.RedirectWithRefresh(w, routes.ProjectConnectionGithubUIDCtx(ctx, ghApp.UID))
	}
}

func HandleGithubInstall(ghCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		installaltion_id := r.URL.Query().Get("installation_id")
		// state := r.URL.Query().Get("source")
		setup_action := r.URL.Query().Get("setup_action")

		installID, err := strconv.ParseInt(installaltion_id, 10, 64)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error parsing installation id")
			render.Error500(w, r)
			return
		}

		ghApp, _ := request.GithubAppFrom(ctx)
		if setup_action != "setup" {
			err = ghCtrl.CompleteInstallation(ctx, ghApp, installID)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error completing github manifest")
				render.Error500(w, r)
				return
			}
		}

		render.RedirectWithRefresh(w, routes.ProjectConnectionGithubUIDCtx(ctx, ghApp.UID))
	}
}
