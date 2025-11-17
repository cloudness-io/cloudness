package volume

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/volume"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog/log"
)

func HandleDelete(volumeCtrl *volume.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		volume, _ := request.VolumeFrom(ctx)

		err := volumeCtrl.SoftDelete(ctx, volume)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error deleting volume")
			render.ToastError(ctx, w, err)
			return
		}

		callerURL, _ := request.HxCallerUrlFrom(ctx)
		render.Redirect(w, callerURL)
	}
}
