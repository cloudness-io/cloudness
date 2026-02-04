package breadcrumb

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/views/components/icons"
	"github.com/cloudness-io/cloudness/app/web/views/dto"
)

func PopulateNavItemKey(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			next.ServeHTTP(w, r.WithContext(
				request.WithNavItem(ctx, &dto.NavItem{Title: key}),
			))
		})
	}
}

func PopulateNavTeam() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tenant, _ := request.TenantFrom(ctx)

			ctx = request.NavItemsReset(ctx)
			next.ServeHTTP(w, r.WithContext(
				request.WithNavItem(ctx, &dto.NavItem{Title: tenant.Name, Icon: icons.TeamIcon, NavURL: routes.TenantCtx(ctx), DropdownActionURL: "/team/list", DropdownIdentifier: "team"}),
			))
		})
	}
}

func PopulateNavProject() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			project, _ := request.ProjectFrom(ctx)
			next.ServeHTTP(w, r.WithContext(
				request.WithNavItem(ctx, &dto.NavItem{
					Title:  project.Name,
					Icon:   icons.ProjectIcon,
					NavURL: routes.ProjectCtx(ctx), DropdownActionURL: routes.TenantCtx(ctx) + "/" + routes.ProjectBase + routes.ProjectNav,
					DropdownIdentifier: "project",
				}),
			))
		})
	}
}
