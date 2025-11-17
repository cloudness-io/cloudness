package source

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/githubapp"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vproject"
	"github.com/cloudness-io/cloudness/app/web/views/components/vsource/vgithubapp"

	"github.com/rs/zerolog/log"
)

func HandleGetGithubApp(ghCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		ghApp, ok := request.GithubAppFrom(ctx)
		target := request.TargetElementFrom(ctx)
		if !ok {
			render.NotFound(w, r)
			return
		}

		if request.HxIndicatorFrom(ctx) && target != routes.TargetMain {
			render.HTML(ctx, w, vgithubapp.Info(tenant, project, ghApp))
			return
		}

		ghApps, err := ghCtrl.List(ctx, tenant.ID, project.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing github apps")
			render.HTML(ctx, w, vgithubapp.Info(tenant, project, ghApp))
			return
		}

		render.Page(ctx, w, vproject.ListGithubApps(tenant, project, GetConfirableSources(), Github, ghApps, ghApp))
	}
}
