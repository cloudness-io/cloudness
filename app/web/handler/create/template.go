package create

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/template"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vcreate"

	"github.com/rs/zerolog/log"
)

func HandleTemplatePreview(templCtrl *template.Controller, templType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		templID, err := request.GetTemplateIDFromPath(r)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Invalid template ID")
			render.ToastError(ctx, w, err)
			return
		}

		templ, err := templCtrl.FindByID(ctx, templID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error finding template")
			render.ToastError(ctx, w, err)
			return
		}

		if templType == "database" {
			render.Page(ctx, w, vcreate.DatabasePreview(templ))
			return
		}

		render.Page(ctx, w, vcreate.TemplatePreview(templ))
	}
}

func HandleTemplateCreate(templCtrl *template.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		templID, err := request.GetTemplateIDFromPath(r)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Invalid template ID")
			render.ToastError(ctx, w, err)
			return
		}

		templ, err := templCtrl.FindByID(ctx, templID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error finding template")
			render.ToastError(ctx, w, err)
			return
		}

		session, _ := request.AuthSessionFrom(ctx)
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		env, _ := request.EnvironmentFrom(ctx)
		err = templCtrl.Create(ctx, session, tenant, project, env, templ)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error creating template")
			render.ToastError(ctx, w, err)
			return
		}

		render.RedirectWithRefresh(w, routes.EnvironmentApplicationCtx(ctx))
	}
}
