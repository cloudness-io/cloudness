package project

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
)

func HandleDelete(projectCtrl *project.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		project, _ := request.ProjectFrom(ctx)

		err := projectCtrl.SoftDelete(ctx, project)
		if err != nil {
			render.ToastError(ctx, w, err)
			return
		}

		render.Redirect(w, routes.TenantCtx(ctx))
	}
}
