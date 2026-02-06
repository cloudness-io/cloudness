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

func HandleGetRestrictions(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)

		restrictions := tenantCtrl.GetRestrctions(ctx, tenant)

		canEdit := canEdit(ctx, tenantCtrl, tenant)
		render.Page(ctx, w, vtenant.Restrictions(tenant, restrictions, canEdit))
	}
}

func HandlePatchRestrictions(tenantCtrl *tenant.Controller) http.HandlerFunc {
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
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}

		render.Page(ctx, w, vtenant.Restrictions(tenant, restrictions, true))
		render.ToastSuccess(ctx, w, "Settings updated successfully")
	}
}
