package inject

import (
	"errors"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/store"

	"github.com/rs/zerolog/log"
)

func InjectApplication(appCtrl *application.Controller) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			applicationUID, err := request.GetApplicationUIDFromPath(r)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error getting application uid from path")
				render.Error500(w, r)
				return
			}

			tenant, _ := request.TenantFrom(ctx)
			project, _ := request.ProjectFrom(ctx)
			environment, _ := request.EnvironmentFrom(ctx)

			application, err := appCtrl.Get(ctx, tenant.ID, project.ID, environment.ID, applicationUID)
			if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
				log.Ctx(ctx).Error().Err(err).Msg("Error fetching application")
				render.Error500(w, r)
				return
			}
			if application == nil {
				render.NotFound(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(
				request.WithApplication(ctx, application),
			))
		})
	}
}
