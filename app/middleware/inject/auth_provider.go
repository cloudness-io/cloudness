package inject

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/auth"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/pages"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

func InjectAuthProvider(authCtrl *auth.Controller) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authProviderStr, err := request.GetAuthProviderFromPath(r)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error getting auth provider from path")
				render.Error500(w, r)
				return
			}

			authProvider := enum.ProviderFromString(authProviderStr)

			authSetting, err := authCtrl.GetByProvider(ctx, authProvider)
			if err != nil || authSetting == nil {
				render.RootWithoutNav(ctx, w, pages.Forbiden(&pages.ForbiddenProps{
					Header:    "Invalid auth provider",
					Subheader: "Please contact administrator",
				}), routes.GetOAuthRedirectUrl(authProvider))
				return
			}

			if !authSetting.Enabled {
				render.RootWithoutNav(ctx, w, pages.Forbiden(&pages.ForbiddenProps{
					Header:    "OAuth provider is disabled",
					Subheader: "Please contact administrator",
				}), routes.GetOAuthRedirectUrl(authProvider))
				return
			}

			next.ServeHTTP(w, r.WithContext(
				request.WithAuthSetting(ctx, authSetting),
			))
		})
	}
}
