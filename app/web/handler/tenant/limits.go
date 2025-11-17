package tenant

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vtenant"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func HandleGetLimits(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)
		canEdit := false

		restrictions := tenantCtrl.GetRestrctions(ctx, tenant)

		if request.IsSuperAdmin(ctx) || (restrictions.AllowAdminToModify && request.IsTeamAdmin(ctx)) {
			canEdit = true
		}

		render.Page(ctx, w, vtenant.Limits(tenant, restrictions, canEdit))
	}
}

func HandlePatchLimits(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)

		in := new(types.TenantRestrictions)
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding request body")
			render.ToastError(ctx, w, err)
			return
		}

		restrictions, err := tenantCtrl.UpdateRestrictions(ctx, tenant, in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error updating tenant restrictions")
			render.ToastError(ctx, w, err)
			return
		}

		render.Page(ctx, w, vtenant.Limits(tenant, restrictions, true))
		render.ToastSuccess(ctx, w, "Settings updated successfully")
	}
}
