package create

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/githubapp"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/handler/source"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vcreate"
	"github.com/cloudness-io/cloudness/app/web/views/dto"

	"github.com/rs/zerolog/log"
)

func HandleListGitOptions(category dto.SourceCategory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		render.Page(ctx, w, vcreate.ListSourceOptions(category, source.GetSources(category)))
	}
}

func HandleListGithubApps(ghAppCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)

		ghApps, err := ghAppCtrl.ListConfigured(ctx, tenant.ID, project.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing github apps")
			render.ToastError(ctx, w, err)
			return
		}

		render.Page(ctx, w, vcreate.ListGithubApps(tenant, project, dto.SourceCategoryGit, source.GetSources(dto.SourceCategoryGit), source.Github, ghApps))
	}
}
