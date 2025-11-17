package variable

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/variable"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog/log"
)

func HandleDelete(appCtrl *application.Controller, varCtrl *variable.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		app, _ := request.ApplicationFrom(ctx)

		variableUID, err := request.GetVariableUIDFromPath(r)
		if err != nil {
			log.Error().Err(err).Msg("Invalid variable uid")
			render.ToastErrorMsg(ctx, w, "Invalid variable uid")
			return
		}

		err = varCtrl.Delete(ctx, app.ID, variableUID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error deleting variable")
			render.ToastError(ctx, w, err)
			return
		}

		appCtrl.NeedsDeployment(ctx, app)

		if err := renderVariablePage(w, r, varCtrl); err == nil {
			render.ToastSuccess(ctx, w, "Variable deleted successfully")
		}
	}
}
