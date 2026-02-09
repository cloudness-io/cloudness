package create

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/gitpublic"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vcreate"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func HandleNewGit(appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		repoURL := r.URL.Query().Get("repoURL")

		//on new page render
		if repoURL == "" {
			render.Page(ctx, w, vcreate.GitPublicPage())
			return
		}

		//on load
		appIn, err := appCtrl.GetGitpublicIn(ctx, repoURL)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting git repository info")
			render.ToastError(ctx, w, err)
			return
		}

		render.Page(ctx, w, vcreate.GitPublicLoadedPage(appIn))
	}
}

func HandleLoadGitPublicRepo(appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		in := new(gitpublic.RepoLoadModel)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding json")
			render.ToastError(ctx, w, err)
			return
		}

		appIn, err := appCtrl.GetGitpublicIn(ctx, in.RepoURL)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting git repository info")
			render.ToastError(ctx, w, err)
			return
		}

		render.Page(ctx, w, vcreate.GitPublicLoadedPage(appIn))
	}
}

func HandleCreateGitPublic(appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)

		in := new(types.ApplicationInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding json")
			render.ToastError(ctx, w, err)
			return
		}

		app, err := appCtrl.Create(ctx, session.Principal.DisplayName, tenant, project, env, in)
		if err != nil {
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}
		ctx = request.WithApplication(ctx, app)

		render.Redirect(w, routes.ApplicationCtx(ctx))
	}
}
