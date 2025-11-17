package inject

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/githubapp"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog/log"
)

func InjectGithubAppSource(ghCtrl *githubapp.Controller) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ghAppUID, err := request.GetSourceUIDFromPath(r)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error getting github app uid from path")
				render.Error500(w, r)
				return
			}

			tenant, _ := request.TenantFrom(ctx)
			project, _ := request.ProjectFrom(ctx)

			ghApp, err := ghCtrl.FindByUID(ctx, tenant.ID, project.ID, ghAppUID)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error fetching github app")
				render.Error500(w, r)
				return
			}
			if ghApp == nil {
				render.NotFound(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(
				request.WithGithubApp(ctx, ghApp),
			))

		})
	}
}
