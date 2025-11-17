package application

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vapplication"
)

func HandleGetTerminal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		application, _ := request.ApplicationFrom(ctx)

		render.Page(ctx, w, vapplication.Terminal(application))
	}
}
