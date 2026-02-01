package tenant

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/dto"
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

		ctx = request.WithNavItem(ctx, &dto.NavItem{Title: "Teams"})
		if request.IsLoginOrRegistrationPage(ctx) {
			w.Header().Set("HX-Purh-Url", "/"+routes.TenantBase)
			render.Page(ctx, w, pages.NewHome(tenants))
			return
		}
		render.RootWithNav(ctx, w, pages.NewHome(tenants), routes.TenantBase)
	}
}
