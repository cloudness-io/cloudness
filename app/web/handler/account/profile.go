package account

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vaccount"
	"github.com/cloudness-io/cloudness/app/web/views/shared"

	"github.com/rs/zerolog/log"
)

func HandleGetProfile(userCtrl *user.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		principal, ok := request.PrincipalFrom(ctx)
		currentURL, _ := request.CurrentFullUrlFrom(ctx)
		if !ok {
			log.Ctx(ctx).Error().Msg("User principal not found")
			render.RootWithoutNav(ctx, w, shared.Maintainance(&shared.MaintainanceProps{
				Header: "You are not logged in",
			}), currentURL)
			return
		}

		user, err := userCtrl.FindUserByEmail(ctx, principal.Email)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting user")
			render.ToastError(ctx, w, err)
			return
		}

		render.Page(ctx, w, vaccount.Profile(user))
	}
}

func HandlePatchProfile(userCtrl *user.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		in := new(user.UpdateInput)
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		principal, _ := request.PrincipalFrom(ctx)

		user, err := userCtrl.Update(ctx, principal.ID, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating profile")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		render.Page(ctx, w, vaccount.Profile(user))
		render.ToastSuccess(ctx, w, "User display name updated successfully")
	}
}
