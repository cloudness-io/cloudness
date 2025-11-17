package inject

import (
	"errors"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/store"

	"github.com/rs/zerolog/log"
)

func InjectTenant(tenantCtrl *tenant.Controller) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			principal, _ := request.PrincipalFrom(ctx)

			tenantUID, err := request.GetTenantUIDFromPath(r)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error getting tenant uid from path")
				render.Error500(w, r)
				return
			}

			tenant, err := tenantCtrl.FindByUID(ctx, tenantUID)
			if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
				log.Ctx(ctx).Error().Err(err).Msg("Error fetching tenant")
				render.Error500(w, r)
				return
			}
			if tenant == nil {
				render.NotFound(w, r)
				return
			}

			membership, err := tenantCtrl.FindMembership(ctx, tenant.ID, principal.ID)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error fetching membership")
				render.Error500(w, r)
				return
			}
			if membership == nil {
				render.NotFound(w, r)
				return
			}

			ctx = request.WithTenant(ctx, tenant)
			ctx = request.WithTenantMembership(ctx, membership)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
