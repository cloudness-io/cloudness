package auth

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/auth"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog/log"
)

func HandleRedirect(authCtrl *auth.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authSettings, _ := request.AuthSettingFrom(ctx)

		redirectUrl, err := authCtrl.GetRedirectUrl(ctx, authSettings)
		if err != nil {
			log.Error().Err(err).Msg("could not get redirect url")
			render.ToastError(ctx, w, err)
			return
		}

		render.RedirectExternal(w, r, redirectUrl)
	}
}
