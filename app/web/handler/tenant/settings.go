package tenant

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vtenant"

	"github.com/rs/zerolog/log"
)

func HandleGetSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)

		render.Page(ctx, w, vtenant.Settings(tenant))
	}
}

func HandlePatchGeneralSettings(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		in := new(tenant.TenantGeneralUpdateModel)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		tenant, _ := request.TenantFrom(ctx)
		tenant, err := tenantCtrl.UpdateGeneral(ctx, tenant, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating tenant")
			render.ToastError(ctx, w, err)
			return
		}

		ctx = request.WithTenant(ctx, tenant)
		render.RootWithNav(ctx, w, vtenant.Settings(tenant), routes.TenantCtx(ctx)+"/"+routes.TenantSettings)
		render.ToastSuccess(ctx, w, "Settings updated successfully")
	}
}
