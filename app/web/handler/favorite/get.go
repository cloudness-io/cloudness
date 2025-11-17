package favorite

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/favorite"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vfavorite"

	"github.com/rs/zerolog/log"
)

func HandleGetFavorite(favCtrl *favorite.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		app, _ := request.ApplicationFrom(ctx)
		session, _ := request.AuthSessionFrom(ctx)

		isFav, err := favCtrl.IsFavorite(ctx, session.Principal.ID, app.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting favorite")
			//TODO reset fav icon
			return
		}

		render.HTML(ctx, w, vfavorite.Symbol(app, isFav))
	}
}
