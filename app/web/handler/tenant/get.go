package tenant

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/pages"

	"github.com/rs/zerolog/log"
)

func HandleGetWithoutTenantUID(tenantCtrl *tenant.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.RedirectWithRefresh(w, routes.TenantBaseURL())
	}
}

func HandleGet(tenantCtrl *tenant.Controller, projectCtrl *project.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		principal, _ := request.PrincipalFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)

		projects, err := projectCtrl.List(ctx, tenant.ID, principal.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error fetching tenant information")
			render.ToastError(ctx, w, err)
			return
		}

		render.RootWithNav(ctx, w, pages.Home(tenant, projects), routes.TenantCtx(ctx))
	}
}
