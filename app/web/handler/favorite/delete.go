package favorite

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/favorite"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vfavorite"

	"github.com/rs/zerolog/log"
)

func HandleDeleteFavorite(favCtrl *favorite.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		app, _ := request.ApplicationFrom(ctx)
		session, _ := request.AuthSessionFrom(ctx)

		err := favCtrl.Delete(ctx, session.Principal.ID, app.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error removing favorite")
			render.ToastError(ctx, w, err)
			return
		}

		render.HTML(ctx, w, vfavorite.Symbol(app, false))
	}
}
