package instance

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/controller/server"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vinstance"

	"github.com/rs/zerolog/log"
)

func HandleGetSettings(instanceCtrl *instance.Controller, serverCtrl *server.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		instance, err := instanceCtrl.Get(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting instance")
			render.ToastError(ctx, w, err)
			return
		}

		server, err := serverCtrl.Get(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting server")
			render.ToastError(ctx, w, err)
			return
		}

		render.Page(ctx, w, vinstance.Settings(instance, server))
	}
}
