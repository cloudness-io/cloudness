package variable

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/variable"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vapplication"

	"github.com/rs/zerolog/log"
)

func renderVariablePage(w http.ResponseWriter, r *http.Request, varCtrl *variable.Controller) error {
	ctx := r.Context()
	application, _ := request.ApplicationFrom(ctx)
	env, _ := request.EnvironmentFrom(ctx)

	vars, err := varCtrl.ListDTO(ctx, env.ID, application.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Error listing environment variables")
		render.ToastError(ctx, w, err)
		return err
	}

	render.Page(ctx, w, vapplication.VariablesNew(application, vars))
	return nil
}
