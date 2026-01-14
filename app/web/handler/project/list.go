package project

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/dto"
	"github.com/cloudness-io/cloudness/app/web/views/shared"

	"github.com/rs/zerolog/log"
)

func HandleListNavigation(projectCtrl *project.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)

		projects, err := projectCtrl.List(ctx, tenant.ID, session.Principal.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing projects")
			render.ToastError(ctx, w, err)
			return
		}

		listItems := make([]*dto.BreadCrumbListItem, 0)
		for _, project := range projects {
			listItems = append(listItems, &dto.BreadCrumbListItem{
				Name: project.Name,
				Link: routes.ProjectCtxUID(ctx, project.UID),
			})
		}

		render.HTML(ctx, w, shared.BreadCrumbDropdownList(listItems))
	}
}
