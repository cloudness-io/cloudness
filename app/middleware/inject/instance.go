package inject

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog/log"
)

func InjectInstance(instanceCtrl *instance.Controller) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			instance, err := instanceCtrl.Get(ctx)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error getting instance settings")
				render.Error500(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(
				request.WithInstanceSettings(ctx, instance),
			))
		})
	}
}
