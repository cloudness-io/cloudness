package tenant

import (
	"encoding/json"
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vtenant"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

func HandleNew() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		render.Page(ctx, w, vtenant.AddTeamPage())
	}
}

func HandleAdd(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)

		in := new(tenant.CreateTenantInput)
		err := json.NewDecoder(r.Body).Decode(in)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error decoding json")
			render.ToastErrorMsg(ctx, w, "Invalid Request Body")
			return
		}

		tenant, err := tenantCtrl.CreateNoAuth(ctx, session.Principal.ID, in, enum.TenantRoleAdmin)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("error creating team")
			render.ToastErrorWithValidation(ctx, w, in, err)
			return
		}
		ctx = request.WithTenant(ctx, tenant)

		render.Redirect(w, routes.TenantCtx(ctx))
	}
}
