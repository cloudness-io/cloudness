package instance

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/controller/server"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vinstance"

	"github.com/rs/zerolog/log"
)

func HandlePatchFQDN(instanceCtrl *instance.Controller, serverCtrl *server.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		in := new(instance.InstanceUpdateFQDNInput)
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}
		server, err := serverCtrl.Get(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting server")
			render.ToastError(ctx, w, err)
			return
		}

		instance, err := instanceCtrl.UpdateFQDN(ctx, server, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating instance settings")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		render.Page(ctx, w, vinstance.Settings(instance, server))
		render.ToastSuccess(ctx, w, "FQDN settings updated successfully")

	}
}

func HandlePatchDNS(instanceCtrl *instance.Controller, serverCtrl *server.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		in := new(instance.InstanceUpdateDNSConfigInput)
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}
		server, err := serverCtrl.Get(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting server")
			render.ToastError(ctx, w, err)
			return
		}

		instance, err := instanceCtrl.UpdateDNSConfig(ctx, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating instance settings")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		render.Page(ctx, w, vinstance.Settings(instance, server))
		render.ToastSuccess(ctx, w, "DNS Settings updated successfully")
	}
}

func HandlePatchScripts(instanceCtrl *instance.Controller, serverCtrl *server.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		in := new(instance.InstanceUpdateScriptsInput)
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}
		server, err := serverCtrl.Get(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting server")
			render.ToastError(ctx, w, err)
			return
		}

		instance, err := instanceCtrl.UpdateScripts(ctx, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating instance settings")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		render.Page(ctx, w, vinstance.Settings(instance, server))
		render.ToastSuccess(ctx, w, "External scripts updated successfully")
	}
}
