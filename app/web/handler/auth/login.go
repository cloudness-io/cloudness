package auth

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/auth"
	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/cookie"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/pages"

	"github.com/rs/zerolog/log"
)

func HandleGet(authCtrl *auth.Controller, instanceCtrl *instance.Controller, cookieName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		instance, _ := request.InstanceSettingsFrom(ctx)
		demoUser, _ := authCtrl.GetDemoUser(ctx, instance)
		enabledAuth, err := authCtrl.ListEnabled(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing auth settings")
			render.Error500(w, r)
			return
		}

		if len(enabledAuth) == 0 {
			render.RootWithoutNav(ctx, w, pages.LoginNoAuth(), routes.Login)
			return
		}

		userSignupEnabled, err := instanceCtrl.IsUserSignupAllowed(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting instance settings")
			render.Error500(w, r)
			return
		}

		cookie.DeleteTokenCookieIfPresent(r, w, cookieName)

		render.RootWithoutNav(ctx, w, pages.LoginPage(userSignupEnabled, enabledAuth, demoUser), routes.Login)
	}
}

func HandleLoginPost(authCtrl *auth.Controller, cookieName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		instance, _ := request.InstanceSettingsFrom(ctx)

		in := new(auth.LoginInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		tokenResponse, err := authCtrl.Login(ctx, instance, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error logging in")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		if cookieName != "" {
			cookie.IncludeTokenCookie(r, w, tokenResponse, cookieName)
		}

		render.RedirectWithRefresh(w, routes.TenantBaseURL())
	}
}
