package server

import (
	"net/http"
	"time"

	"github.com/cloudness-io/cloudness/app/controller/server"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vserver"

	"github.com/rs/zerolog/log"
)

func HandleListCertificates(serverCtrl *server.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		time.Sleep(1 * time.Second)
		certs, err := serverCtrl.ListCertificates(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing certificates")
			render.ToastErrorMsg(ctx, w, "Error listing certificates")
			return
		}

		render.HTML(ctx, w, vserver.ListCertificates(certs))
	}
}
