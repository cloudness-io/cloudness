package application

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vapplication"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

func HandleGetMetrics(appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		application, _ := request.ApplicationFrom(ctx)

		metricsSpan, _ := enum.MetricsSpan(request.GetMetricsSpanFromPath(r)).Sanitize()

		metrics, err := appCtrl.GetMetrics(ctx, application, metricsSpan)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting metrics")
			render.ToastError(ctx, w, err)
			return
		}

		render.Page(ctx, w, vapplication.Metrics(application, metrics, metricsSpan))
	}
}
