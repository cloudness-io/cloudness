package variable

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/variable"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog/log"
)

func HandlePost(appCtrl *application.Controller, varCtrl *variable.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		app, _ := request.ApplicationFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)

		in := new(variable.AddVariableInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Error().Err(err).Msg("Invalid request body")
			render.ToastErrorMsg(ctx, w, "Invalid request body")
			return
		}

		err := varCtrl.Add(ctx, env.ID, app.ID, in)
		if err != nil {
			log.Error().Err(err).Msg("Error updating variables")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		appCtrl.NeedsDeployment(ctx, app)

		if err := renderVariablePage(w, r, varCtrl); err == nil {
			render.ToastSuccess(ctx, w, "Variables added successfully")
		}
	}
}
