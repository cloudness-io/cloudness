package logs

import (
	"errors"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/logs"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vlog"
	"github.com/cloudness-io/cloudness/store"

	"github.com/rs/zerolog/log"
)

func HandleGetLogs(logsCtrl *logs.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		deployment, _ := request.DeploymentFrom(ctx)

		lines, err := logsCtrl.Find(ctx, deployment.ID)
		if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
			log.Err(err).Msg("error fetching logs")
			render.ToastError(ctx, w, err)
			return
		}

		render.HTML(ctx, w, vlog.LogViewer(lines))
	}
}
