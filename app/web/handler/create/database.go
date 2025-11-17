package create

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/template"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vcreate"

	"github.com/rs/zerolog/log"
)

func HandleGetDatabaseView(templCtrl *template.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tmpls, err := templCtrl.ListDatabase(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing database templates")
			render.ToastError(ctx, w, err)
		}

		render.Page(ctx, w, vcreate.DatabaseView(tmpls))
	}
}
