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

func HandlePatch(appCtrl *application.Controller, varCtrl *variable.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		env, _ := request.EnvironmentFrom(ctx)
		app, _ := request.ApplicationFrom(ctx)
		variableUID, err := request.GetVariableUIDFromPath(r)
		if err != nil {
			log.Error().Err(err).Msg("Invalid variable uid")
			render.ToastErrorMsg(ctx, w, "Invalid variable uid")
			return
		}

		in := new(variable.UpdateVariableInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Error().Err(err).Msg("Invalid request body")
			render.ToastErrorMsg(ctx, w, "Invalid request body")
			return
		}

		if err := varCtrl.Update(ctx, env.ID, app.ID, variableUID, in); err != nil {
			log.Error().Err(err).Msg("error updating variable")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		appCtrl.NeedsDeployment(ctx, app)

		if err := renderVariablePage(w, r, varCtrl); err == nil {
			render.ToastSuccess(ctx, w, "Variable updated successfully")
		}
	}
}

func HandleGenerate(appCtrl *application.Controller, varCtrl *variable.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		app, _ := request.ApplicationFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)
		variableUID, err := request.GetVariableUIDFromPath(r)
		if err != nil {
			log.Error().Err(err).Msg("Invalid variable uid")
			render.ToastErrorMsg(ctx, w, "Invalid variable uid")
			return
		}

		in := new(variable.GenerateVariableInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Error().Err(err).Msg("Invalid request body")
			render.ToastErrorMsg(ctx, w, "Invalid request body")
			return
		}

		if err := varCtrl.UpdateGenerate(ctx, env.ID, app.ID, variableUID, in); err != nil {
			log.Error().Err(err).Msg("error generating variable")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		appCtrl.NeedsDeployment(ctx, app)

		if err := renderVariablePage(w, r, varCtrl); err == nil {
			render.ToastSuccess(ctx, w, "Variable generated successfully")
		}
	}
}
