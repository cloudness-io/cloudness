package application

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/githubapp"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog/log"
)

func HandleAddTCPProxy(appCtrl *application.Controller, ghAppCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)
		app, _ := request.ApplicationFrom(ctx)

		in := new(application.TCPProxyInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		_, err := appCtrl.AddTCPProxy(ctx, session, tenant, project, env, app, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error adding tcp proxy")
			render.ToastError(ctx, w, err)
			return
		}
		if err := renderSettings(ctx, w, appCtrl, ghAppCtrl); err == nil {
			render.ToastSuccess(ctx, w, "TCP proxy added successfully")
		}
	}
}

func HandleDeleteTCPProxy(appCtrl *application.Controller, ghAppCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)
		app, _ := request.ApplicationFrom(ctx)

		_, err := appCtrl.DeleteTCPProxy(ctx, session, tenant, project, env, app)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error deleting tcp proxy")
			render.ToastError(ctx, w, err)
			return
		}
		if err := renderSettings(ctx, w, appCtrl, ghAppCtrl); err == nil {
			render.ToastSuccess(ctx, w, "TCP Proxy deleted successfully")
		}
	}
}
