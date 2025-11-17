package application

import (
	"net/http"
	"strconv"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vapplication"

	"github.com/rs/zerolog/log"
)

func HandleDeleteView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		application, _ := request.ApplicationFrom(ctx)

		render.Page(ctx, w, vapplication.Delete(application))
	}
}

func HandleDeleteApplication(appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		app, _ := request.ApplicationFrom(ctx)

		opts := new(application.AppDeleteOption)
		volumeStr := r.URL.Query().Get("volume")
		if volumeStr != "" {
			volume, err := strconv.ParseBool(volumeStr)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error parsing volume query param")
				render.ToastError(ctx, w, err)
				return
			}
			opts.Volume = volume
		}
		err := appCtrl.SoftDelete(ctx, app, opts)
		if err != nil {
			log.Err(err).Msg("error deleting application")
			render.ToastError(ctx, w, err)
			return
		}

		render.Redirect(w, routes.EnvironmentCtx(ctx)+routes.EnvironmentApplication+routes.TargetMainQuery)
	}
}
