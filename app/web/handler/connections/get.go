package connections

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/githubapp"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vsource/vgithubapp"
)

func HandleGetGithubApp(ghCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)
		ghApp, _ := request.GithubAppFrom(ctx)

		render.Page(ctx, w, vgithubapp.GHAppInfoPage(tenant, project, ghApp))
	}
}
