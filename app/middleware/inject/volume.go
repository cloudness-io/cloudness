package inject

import (
	"errors"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/volume"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/store"

	"github.com/rs/zerolog/log"
)

func InjectVolume(volumeCtrl *volume.Controller) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			volumeUID, err := request.GetVolumeUIDFromPath(r)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error getting volume uid from path")
				render.Error500(w, r)
				return
			}

			tenant, _ := request.TenantFrom(ctx)
			project, _ := request.ProjectFrom(ctx)
			environment, _ := request.EnvironmentFrom(ctx)

			volume, err := volumeCtrl.Get(ctx, tenant.ID, project.ID, environment.ID, volumeUID)
			if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
				log.Ctx(ctx).Error().Err(err).Msg("Error fetching volume")
				render.Error500(w, r)
				return
			}
			if volume == nil {
				render.NotFound(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(
				request.WithVolume(ctx, volume),
			))
		})
	}
}
