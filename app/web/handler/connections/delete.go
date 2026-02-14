package connections

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/githubapp"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog/log"
)

func HandleDeleteGithubApp(ghCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ghApp, _ := request.GithubAppFrom(ctx)

		err := ghCtrl.Delete(ctx, ghApp)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error deleting github app")
			render.ToastError(ctx, w, err)
			return
		}

		render.Redirect(w, routes.ProjectCtx(ctx)+"/"+routes.ProjectConnections)
	}
}
