package create

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vcreate"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func HandleGetRegistryView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		render.Page(ctx, w, vcreate.RegistryView())
	}
}

func HandleCreateWithRegistry(appCtrl *application.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		in := new(types.RegistryInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)

		app, err := appCtrl.CreateRegistry(ctx, session.Principal.DisplayName, tenant, project, env, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error creating application")
			render.ToastError(ctx, w, err)
			return
		}
		ctx = request.WithApplication(ctx, app)

		render.Redirect(w, routes.ApplicationCtx(ctx))
	}
}
