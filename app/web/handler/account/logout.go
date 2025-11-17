package account

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/auth"
	"github.com/cloudness-io/cloudness/app/cookie"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog/log"
)

func HandleLogout(authCtrl *auth.Controller, cookieName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		err := authCtrl.Logout(ctx, session)

		// best effort elete cookie even in case of errors, to avoid clients being unable to remove the cookie.
		// WARNING: It could be that the cookie is removed even though the token is still there in the DB.
		// However, we have APIs to list and delete session tokens, and expiry time is usually short.
		cookie.DeleteTokenCookieIfPresent(r, w, cookieName)

		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error logging out")
			return
		}
		// if err != nil {
		// 	// TODO: handle error ??
		// 	log.Err(err).Msg("Error logging out")
		// 	return
		// }

		// loginUrl, err := getLoginUrl(oauthCtrl)
		// if err != nil {
		// 	log.Ctx(ctx).Error().Err(err).Msg("Error getting login url")
		// 	render.Error500(w, r)
		// 	return
		// }

		render.RedirectWithRefresh(w, "/")
	}
}
