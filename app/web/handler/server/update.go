package server

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/server"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vinstance"

	"github.com/rs/zerolog/log"
)

func HandlePatchGeneral(serverCtrl *server.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		in := new(server.ServerGeneralUpdateModel)
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		server, err := serverCtrl.UpdateGeneral(ctx, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating server")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		render.Page(ctx, w, vinstance.Server(server))
		render.ToastSuccess(ctx, w, "Server settings updated successfully")
	}
}

func HandlePatchNetwork(serverCtrl *server.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		in := new(server.ServerNetworkUpdateModel)
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		server, err := serverCtrl.UpdateNetwork(ctx, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating server")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		render.Page(ctx, w, vinstance.Server(server))
		render.ToastSuccess(ctx, w, "Server network settings updated successfully")
	}
}

func HandlePatchBuilder(serverCtrl *server.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		in := new(server.ServerBuilderUpdateModel)
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		server, err := serverCtrl.UpdateBuilder(ctx, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating server")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		render.Page(ctx, w, vinstance.Server(server))
		render.ToastSuccess(ctx, w, "Server builder settings updated successfully")
	}
}

func HandlePatchLimits(serverCtrl *server.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		in := new(server.ServerLimitsUpdateModel)
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		server, err := serverCtrl.UpdateLimits(ctx, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating server")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		render.Page(ctx, w, vinstance.Server(server))
		render.ToastSuccess(ctx, w, "Server Limits updated successfully")
	}
}
