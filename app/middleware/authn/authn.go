package authn

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cloudness-io/cloudness/app/auth/authn"
	"github.com/cloudness-io/cloudness/app/controller/auth"
	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// AttemptWeb returns an http.HandlerFunc middleware that authenticates
// the http.Request if authentication html payload is available.
// Otherwise, user will be redirected to login
func AttemptWeb(authenticator authn.Authenticator, instanceCtrl *instance.Controller, authCtrl *auth.Controller, cookieName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := hlog.FromRequest(r)

			session, err := authenticator.Authenticate(r)
			if err != nil {
				log.Error().Err(err).Msg("authentication failed")
				if errors.Is(err, authn.ErrNoAuthData) || strings.HasPrefix(err.Error(), "JWT:") {
					render.RedirectExternal(w, r, routes.Login)
					return
				} else {
					render.Error500(w, r)
					return
				}
			}

			// Update the logging context and inject principal in context
			log.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.
					Str("principal_name", session.Principal.DisplayName).
					Str("principal_type", string(session.Principal.Type))
			})

			next.ServeHTTP(w, r.WithContext(
				request.WithAuthSession(ctx, session),
			))
		})
	}
}
