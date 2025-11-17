package account

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vaccount"
	"github.com/cloudness-io/cloudness/app/web/views/shared"

	"github.com/rs/zerolog/log"
)

func HandleGetSession(userCtrl *user.Controller) http.HandlerFunc {
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

		render.Page(ctx, w, vaccount.Session(user))
	}
}
