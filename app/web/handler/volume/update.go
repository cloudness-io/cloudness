package volume

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/volume"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func HandleUpdateAttached(appCtrl *application.Controller, volumeCtrl *volume.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		application, _ := request.ApplicationFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)
		session, _ := request.AuthSessionFrom(ctx)
		volume, _ := request.VolumeFrom(ctx)

		in := new(types.VolumeCreateInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		volume.FromInput(in)

		_, err := appCtrl.UpdateVolume(ctx, session, tenant, project, env, application, volume)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating volume")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}
		err = renderVolumesList(ctx, w, application, volumeCtrl)
		if err == nil {
			render.ToastSuccess(ctx, w, "Volume updated successfully")
		}
	}
}

func HandleUpdateDetach(appCtrl *application.Controller, volumeCtrl *volume.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		application, _ := request.ApplicationFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)
		session, _ := request.AuthSessionFrom(ctx)
		volume, _ := request.VolumeFrom(ctx)

		_, err := appCtrl.DetachVolume(ctx, session, tenant, project, env, application, volume)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating volume")
			render.ToastError(ctx, w, err)
			return
		}
		err = renderVolumesList(ctx, w, application, volumeCtrl)
		if err == nil {
			render.ToastSuccess(ctx, w, "Volume detached successfully")
		}
	}
}
