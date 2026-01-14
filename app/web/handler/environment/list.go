package environment

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/environment"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vproject"
	"github.com/cloudness-io/cloudness/app/web/views/dto"
	"github.com/cloudness-io/cloudness/app/web/views/shared"

	"github.com/rs/zerolog/log"
)

func HandleList(envCtrl *environment.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)

		envs, err := envCtrl.List(ctx, tenant.ID, project.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing environments")
			render.ToastError(ctx, w, err)
			return
		}

		render.Page(ctx, w, vproject.ListEnvs(project, envs, nil))
	}
}

func HandleListNavigation(envCtrl *environment.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)

		envs, err := envCtrl.List(ctx, tenant.ID, project.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing environments")
			render.ToastError(ctx, w, err)
			return
		}

		listItems := make([]*dto.BreadCrumbListItem, 0)
		for _, env := range envs {
			listItems = append(listItems, &dto.BreadCrumbListItem{
				Name: env.Name,
				Link: routes.EnvironmentCtxUID(ctx, env.UID) + routes.EnvironmentApplication + routes.TargetMainQuery,
			})
		}

		render.HTML(ctx, w, shared.BreadCrumbDropdownList(listItems))
	}
}
