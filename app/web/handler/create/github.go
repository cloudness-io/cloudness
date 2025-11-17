package create

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vcreate"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func HandleGetGithubView(appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ghApp, _ := request.GithubAppFrom(ctx)

		in := appCtrl.GetDefaultGitIn()
		render.Page(ctx, w, vcreate.Github(ghApp, in))
	}
}

func HandleCreateGithub(appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)
		ghApp, _ := request.GithubAppFrom(ctx)

		in := new(types.ApplicationInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding json")
			render.ToastError(ctx, w, err)
			return
		}

		app, err := appCtrl.CreateGithub(ctx, session.Principal.DisplayName, tenant, project, env, nil, ghApp, in)
		if err != nil {
			render.ToastError(ctx, w, err)
			return
		}
		ctx = request.WithApplication(ctx, app)

		render.Redirect(w, routes.ApplicationCtx(ctx))
	}
}
