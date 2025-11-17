package auth

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/auth"
	"github.com/cloudness-io/cloudness/app/cookie"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/usererror"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/shared"

	"github.com/rs/zerolog/log"
)

func HandleCallback(authCtrl *auth.Controller, cookieName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authSettings, _ := request.AuthSettingFrom(ctx)

		code := r.URL.Query().Get("code")

		tokenResponse, err := authCtrl.Callback(ctx, authSettings, code)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error logging in")
			render.RootWithoutNav(ctx, w, shared.Maintainance(&shared.MaintainanceProps{
				Header:    usererror.TranslateErrMsg(ctx, err),
				Subheader: "Please try again later. If the issue persists, please contact support.",
			}), routes.GetOAuthCallbackUrl(authSettings.Provider))
			return
		}

		if cookieName != "" {
			cookie.IncludeTokenCookie(r, w, tokenResponse, cookieName)
		}

		render.RedirectWithRefresh(w, "/")

	}
}
