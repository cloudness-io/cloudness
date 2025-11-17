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

func HandleUpdatePrivateDomain(appCtrl *application.Controller, ghAppCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)
		app, _ := request.ApplicationFrom(ctx)

		in := new(application.PrivateNetworkInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		_, err := appCtrl.UpdatePrivateNetwork(ctx, session, tenant, project, env, app, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating private network")
			render.ToastError(ctx, w, err)
			return
		}
		if err := renderSettings(ctx, w, appCtrl, ghAppCtrl); err == nil {
			render.ToastSuccess(ctx, w, "Private networking updated successfully")
		}

	}
}

func HandleGenerateDomain(appCtrl *application.Controller, ghAppCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)
		app, _ := request.ApplicationFrom(ctx)

		_, err := appCtrl.GenerateDomain(ctx, session, tenant, project, env, app)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error generating domain")
			render.ToastError(ctx, w, err)
			return
		}
		if err := renderSettings(ctx, w, appCtrl, ghAppCtrl); err == nil {
			render.ToastSuccess(ctx, w, "Domain generated successfully")
		}
	}
}

func HandleUpdateDomain(appCtrl *application.Controller, ghAppCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)
		app, _ := request.ApplicationFrom(ctx)

		in := new(application.UpdateDomainInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		_, err := appCtrl.UpdateDomain(ctx, session, tenant, project, env, app, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating domain")
			render.ToastError(ctx, w, err)
			return
		}
		if err := renderSettings(ctx, w, appCtrl, ghAppCtrl); err == nil {
			render.ToastSuccess(ctx, w, "Domain updated successfully")
		}
	}
}

func HandleDeleteDomain(appCtrl *application.Controller, ghAppCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)
		app, _ := request.ApplicationFrom(ctx)

		_, err := appCtrl.DeleteDomain(ctx, session, tenant, project, env, app)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error deleting domain")
			render.ToastError(ctx, w, err)
			return
		}
		if err := renderSettings(ctx, w, appCtrl, ghAppCtrl); err == nil {
			render.ToastSuccess(ctx, w, "Domain deleted successfully")
		}
	}
}
