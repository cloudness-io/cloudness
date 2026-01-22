package create

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/template"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vcreate"
	"github.com/rs/zerolog/log"
)

func HandleGetOnelickTemplate(templCtrl *template.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		templates, err := templCtrl.ListTemplates(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing templates")
			return
		}

		render.Page(ctx, w, vcreate.OneClick(templates))
	}
}
