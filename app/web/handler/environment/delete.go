package environment

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/environment"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"
)

func HandleDelete(envCtrl *environment.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		env, _ := request.EnvironmentFrom(ctx)

		err := envCtrl.SoftDelete(ctx, env)
		if err != nil {
			render.ToastError(ctx, w, err)
			return
		}

		render.Redirect(w, routes.ProjectEnvironment)
	}
}
