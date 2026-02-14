package instance

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vinstance"

	"github.com/rs/zerolog/log"
)

func HandleGetRegistry(instanceCtrl *instance.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		instance, err := instanceCtrl.Get(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting instance")
			render.ToastError(ctx, w, err)
			return
		}
		render.Page(ctx, w, vinstance.Registry(instance))
	}
}

func HandlePatchRegistry(instanceCtrl *instance.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		in := new(instance.InstanceRegistryUpdateInput)
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		instance, err := instanceCtrl.UpdateRegistry(ctx, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating instance registry settings")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		render.Page(ctx, w, vinstance.Registry(instance))
		render.ToastSuccess(ctx, w, "Registry settings updated successfully")
	}
}
