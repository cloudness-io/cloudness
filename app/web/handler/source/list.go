package source

import (
	"net/http"

	"github.com/cloudness-io/cloudness/app/controller/githubapp"
	"github.com/cloudness-io/cloudness/app/controller/gitpublic"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/components/vgit"
	"github.com/cloudness-io/cloudness/app/web/views/components/vproject"
	"github.com/cloudness-io/cloudness/helpers"

	"github.com/rs/zerolog/log"
)

func HandleListConfigurableSource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		project, _ := request.ProjectFrom(ctx)

		render.Page(ctx, w, vproject.ListOptions(project, GetConfirableSources()))
	}
}

func HandleListGithubApps(ghCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenant, _ := request.TenantFrom(ctx)
		project, _ := request.ProjectFrom(ctx)

		ghApps, err := ghCtrl.List(ctx, tenant.ID, project.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing github apps")
			render.ToastError(ctx, w, err)
			return
		}

		render.Page(ctx, w, vproject.ListGithubApps(tenant, project, GetConfirableSources(), Github, ghApps, nil))
	}
}

func HandleListGithubRepos(ghCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		selected := r.URL.Query().Get("selected")
		ghApp, _ := request.GithubAppFrom(ctx)

		repos, err := ghCtrl.ListRepos(ctx, ghApp)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing repos")
			render.ToastError(ctx, w, err)
			return
		}

		render.HTML(ctx, w, vgit.ListRepos(repos, selected))
	}
}

func HandleListGithubBranches(ghCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ghApp, _ := request.GithubAppFrom(ctx)
		selected := r.URL.Query().Get("selected")
		fullName := r.URL.Query().Get("repo")

		owner, repo, err := helpers.SplitGitRepoFullname(fullName)
		if err != nil {
			render.ToastError(ctx, w, err)
			return
		}

		branches, err := ghCtrl.ListBranches(ctx, ghApp, owner, repo)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing branches")
			render.ToastError(ctx, w, err)
			return
		}
		if selected == "" && len(branches) > 0 {
			selected = branches[0].Name
		}

		render.HTML(ctx, w, vgit.ListBranches(branches, selected))
	}
}

func HandleListGitpublicBranches(gitPublicCtrl *gitpublic.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		selected := r.URL.Query().Get("branch")
		repoURL := r.URL.Query().Get("repoURL")

		branches, err := gitPublicCtrl.ListBranches(ctx, repoURL)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error listing branches")
			render.ToastError(ctx, w, err)
			return
		}
		if selected == "" && len(branches) > 0 {
			selected = branches[0].Name
		}

		render.HTML(ctx, w, vgit.ListBranches(branches, selected))
	}
}
