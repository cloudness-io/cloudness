package inject

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/deployment"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog/log"
)

func InjectDeployment(deploymentCtrl *deployment.Controller) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			deploymentUID, err := request.GetDeploymentUIDFromPath(r)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error getting deployment uid from path")
				render.Error500(w, r)
				return
			}

			application, _ := request.ApplicationFrom(ctx)

			deployment, err := deploymentCtrl.Get(ctx, application.ID, deploymentUID)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error fetching deployment")
				render.Error500(w, r)
				return
			}
			if deployment == nil {
				render.NotFound(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(
				request.WithDeployment(ctx, deployment),
			))
		})
	}
}
