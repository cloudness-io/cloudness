package tenant

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vtenant"
	"github.com/cloudness-io/cloudness/app/web/views/pages"

	"github.com/rs/zerolog/log"
)

func HandleList(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		principal, _ := request.PrincipalFrom(ctx)

		tenants, err := tenantCtrl.ListMembership(ctx, principal.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting tenants")
			render.ToastError(ctx, w, err)
			return
		}

		if request.IsLoginOrRegistrationPage(ctx) {
			w.Header().Set("HX-Purh-Url", routes.TenantBaseURL())
			render.Page(ctx, w, pages.Home(tenants))
			return
		}

		render.RootWithNav(ctx, w, pages.Home(tenants), routes.TenantBaseURL())
	}
}

func HandleListNavigation(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		principal, _ := request.PrincipalFrom(ctx)
		selectedUID, _ := request.GetSelectedUIDFromPath(r)

		projects, err := tenantCtrl.ListMembership(ctx, principal.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing projects")
			render.ToastError(ctx, w, err)
			return
		}

		render.HTML(ctx, w, vtenant.Dropdown(projects, selectedUID))
	}
}
