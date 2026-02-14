package project

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vproject"

	"github.com/rs/zerolog/log"
)

func HandleListNavigation(projectCtrl *project.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		selectedUID, _ := request.GetSelectedUIDFromPath(r)

		projects, err := projectCtrl.List(ctx, tenant.ID, session.Principal.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing projects")
			render.ToastError(ctx, w, err)
			return
		}

		render.HTML(ctx, w, vproject.Dropdown(projects, tenant, selectedUID))
	}
}
