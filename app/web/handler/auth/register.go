package auth

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/auth"
	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/cookie"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/pages"

	"github.com/rs/zerolog/log"
)

func HandleGetRegister(instanceCtrl *instance.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		useSignupEnabled, err := instanceCtrl.IsUserSignupAllowed(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting instance settings")
			render.Error500(w, r)
			return
		}

		if !useSignupEnabled {
			render.RootWithoutNav(ctx, w, pages.RegistrationDisabled(), routes.Register)
			return
		}

		render.RootWithoutNav(ctx, w, pages.RegistrationPage(), routes.Register)
	}
}

func HandleRegister(instanceCtrl *instance.Controller, authCtrl *auth.Controller, cookieName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		useSignupEnabled, err := instanceCtrl.IsUserSignupAllowed(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting instance settings")
			render.Error500(w, r)
			return
		}

		if !useSignupEnabled {
			render.RootWithoutNav(ctx, w, pages.RegistrationDisabled(), routes.Register)
			return
		}
		in := new(auth.RegisterUserInput)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		tokenResponse, err := authCtrl.Register(ctx, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error registering user")
			render.ToastError(ctx, w, err)
			return
		}

		if cookieName != "" {
			cookie.IncludeTokenCookie(r, w, tokenResponse, cookieName)
		}

		render.RedirectWithRefresh(w, "/")
	}
}
