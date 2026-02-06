package tenant

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vtenant"

	"github.com/rs/zerolog/log"
)

func HandleGetSettings(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)

		canEdit := canEdit(ctx, tenantCtrl, tenant)
		render.Page(ctx, w, vtenant.Settings(tenant, canEdit))
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
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		canEdit := canEdit(ctx, tenantCtrl, tenant)

		ctx = request.WithTenant(ctx, tenant)
		render.Page(ctx, w, vtenant.Settings(tenant, canEdit))
		render.ToastSuccess(ctx, w, "Settings updated successfully")
	}
}
