package deployment

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/deployment"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vapplication"
)

func HandleGetDeployment(deploymentCtrl *deployment.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		application, _ := request.ApplicationFrom(ctx)
		deployment, _ := request.DeploymentFrom(ctx)

		render.Page(ctx, w, vapplication.DeploymentInfo(application, deployment))
	}
}
