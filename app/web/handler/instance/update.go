package instance

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/controller/server"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vinstance"

	"github.com/rs/zerolog/log"
)

func HandlePatchInstanceSettings(instanceCtrl *instance.Controller, serverCtrl *server.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		in := new(instance.InstanceUpdateInput)
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}
		server, err := serverCtrl.Get(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting server")
			render.ToastError(ctx, w, err)
			return
		}

		instance, err := instanceCtrl.Update(ctx, server, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating instance settings")
			render.ToastError(ctx, w, err)
			return
		}

		render.Page(ctx, w, vinstance.Settings(instance, server))
		render.ToastSuccess(ctx, w, "Instance settings updated successfully")
	}
}
