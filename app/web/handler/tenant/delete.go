package tenant

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vtenant"

	"github.com/rs/zerolog/log"
)

func HandleGetDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)

		render.Page(ctx, w, vtenant.Delete(tenant))
	}
}

func HandleDeleteTeam(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)

		err := tenantCtrl.SoftDelete(ctx, tenant)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error deleting tenant")
			render.ToastError(ctx, w, err)
			return
		}

		render.RedirectWithRefresh(w, "/")
	}
}
