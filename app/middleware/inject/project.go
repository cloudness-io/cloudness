package inject

import (
	"errors"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/store"

	"github.com/rs/zerolog/log"
)

func InjectProject(projectCtrl *project.Controller) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tenant, _ := request.TenantFrom(ctx)
			principal, _ := request.PrincipalFrom(ctx)

			projectUID, err := request.GetProjectUIDFromPath(r)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error getting project uid from path")
				render.Error500(w, r)
				return
			}

			project, err := projectCtrl.FindByUID(ctx, tenant.ID, projectUID)
			if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
				log.Ctx(ctx).Error().Err(err).Msg("Error fetching deployment")
				render.Error500(w, r)
				return
			}
			if project == nil {
				render.NotFound(w, r)
				return
			}

			if !request.IsTeamAdmin(ctx) {
				membership, err := projectCtrl.FindMembership(ctx, tenant.ID, project.ID, principal.ID)
				if err != nil {
					log.Ctx(ctx).Error().Err(err).Msg("Error fetching membership")
					render.Error500(w, r)
					return
				}
				if membership == nil {
					render.NotFound(w, r)
					return
				} else {
					ctx = request.WithProjectMembership(ctx, membership)
				}
			}

			ctx = request.WithProject(ctx, project)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
