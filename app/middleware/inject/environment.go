package inject

import (
	"errors"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/environment"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/store"

	"github.com/rs/zerolog/log"
)

func InjectEnvironment(envCtrl *environment.Controller) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			project, _ := request.ProjectFrom(ctx)

			envUID, err := request.GetEnvironmentUIDFromPath(r)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error getting env uid from path")
				render.Error500(w, r)
				return
			}

			environment, err := envCtrl.Get(ctx, project.ID, envUID)
			if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
				log.Ctx(ctx).Error().Err(err).Msg("Error fetching environment")
				render.Error500(w, r)
				return
			}
			if environment == nil {
				render.NotFound(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(
				request.WithEnvironment(ctx, environment),
			))
		})
	}
}
