package volume

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/volume"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vapplication"
	"github.com/cloudness-io/cloudness/app/web/views/components/vvolume"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

func HandleListVolume(volumeCtrl *volume.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		application, _ := request.ApplicationFrom(ctx)

		volumes, err := volumeCtrl.ListForApp(ctx, application)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing volumes")
			render.ToastError(ctx, w, err)
			return
		}

		render.Page(ctx, w, vapplication.VolumeList(application, volumes))
	}
}

func HandleListUnattached(volumeCtrl *volume.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)

		volumes, err := volumeCtrl.ListUnattached(ctx, tenant.ID, project.ID, env.ID, nil)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing unattached volumes")
			render.ToastError(ctx, w, err)
			return
		}

		render.HTML(ctx, w, vvolume.ListUnattached(volumes, enum.VolumeFormActionDelete))
	}
}

func HandleListAttachable(volumeCtrl *volume.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)
		app, _ := request.ApplicationFrom(ctx)

		volumes, err := volumeCtrl.ListUnattached(ctx, tenant.ID, project.ID, env.ID, &app.ServerID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing unattached volumes")
			render.ToastError(ctx, w, err)
			return
		}

		render.HTML(ctx, w, vvolume.ListUnattached(volumes, enum.VolumeFormActionAttach))
	}
}
