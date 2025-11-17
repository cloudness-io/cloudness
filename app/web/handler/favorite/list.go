package favorite

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/favorite"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vfavorite"

	"github.com/rs/zerolog/log"
)

func HandleListFavorites(favCtrl *favorite.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)
		session, _ := request.AuthSessionFrom(ctx)

		favs, err := favCtrl.List(ctx, session.Principal.ID, tenant.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing favorites")
			render.ToastErrorMsg(ctx, w, "Error listing favorites")
			return
		}

		render.HTML(ctx, w, vfavorite.List(favs))
	}
}
