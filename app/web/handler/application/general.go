package application

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog/log"
)

func HandleUpdateIcon(appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		in := new(application.UpdateIconInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		app, _ := request.ApplicationFrom(ctx)

		_, err := appCtrl.UpdateIcon(ctx, app, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating icon")
			render.ToastError(ctx, w, err)
			return
		}

		render.Redirect(w, routes.ApplicationCtx(ctx))
	}
}
