package volume

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/volume"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vapplication"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func HandleCreate(appCtrl *application.Controller, volumeCtrl *volume.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		application, _ := request.ApplicationFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)
		session, _ := request.AuthSessionFrom(ctx)

		in := new(types.VolumeCreateInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		_, err := appCtrl.AddVolume(ctx, session, tenant, project, env, application, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error creating volume")
			render.ToastError(ctx, w, err)
			return
		}
		err = renderVolumesList(ctx, w, application, volumeCtrl)
		if err == nil {
			render.ToastSuccess(ctx, w, "Volume created successfully")
		}
	}
}

func renderVolumesList(ctx context.Context, w http.ResponseWriter, application *types.Application, volumeCtrl *volume.Controller) error {
	volumes, err := volumeCtrl.ListForApp(ctx, application)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Error listing volumes")
		render.ToastError(ctx, w, err)
		return err
	}

	render.Page(ctx, w, vapplication.VolumeList(application, volumes))
	return nil
}
