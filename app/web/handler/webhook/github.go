package webhook

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudness-io/cloudness/app/controller/githubapp"
	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/render"

	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/errors"
)

type githubAppState struct {
	tenantUID  int64
	projectUID int64
	ghAppUID   int64
}

func parseGithubAppState(state string) (*githubAppState, error) {
	segments := strings.Split(state, "-")
	if len(segments) != 3 {
		return nil, errors.NewBadRequest("Invalid redirect state")
	}
	tenantUID, err := strconv.ParseInt(segments[0], 10, 64)
	if err != nil {
		return nil, errors.NewBadRequest("Invalid redirect state")
	}
	projectUID, err := strconv.ParseInt(segments[1], 10, 64)
	if err != nil {
		return nil, errors.NewBadRequest("Invalid redirect state")
	}
	ghAppUID, err := strconv.ParseInt(segments[2], 10, 64)
	if err != nil {
		return nil, errors.NewBadRequest("Invalid redirect state")
	}
	return &githubAppState{tenantUID: tenantUID, projectUID: projectUID, ghAppUID: ghAppUID}, nil
}

func HandleGithubRedirect(tenantCtrl *tenant.Controller, projectCtrl *project.Controller, ghCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		ghAppState, err := parseGithubAppState(state)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error parsing github app uid from state")
			render.Error500(w, r)
			return
		}

		tenant, err := tenantCtrl.FindByUID(ctx, ghAppState.tenantUID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error finding tenant")
			render.Error500(w, r)
			return
		}

		project, err := projectCtrl.FindByUID(ctx, tenant.ID, ghAppState.projectUID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error finding tenant")
			render.Error500(w, r)
			return
		}

		ghApp, err := ghCtrl.FindByUID(ctx, tenant.ID, project.ID, ghAppState.ghAppUID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting github app")
			render.Error500(w, r)
			return
		}

		if ghApp == nil {
			log.Ctx(ctx).Error().Err(err).Msg("Github app not found")
			render.Error500(w, r)
			return
		}

		err = ghCtrl.CompleteManifest(ctx, ghApp, code)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error completing github manifest")
			render.Error500(w, r)
			return
		}

		ctx = request.WithTenant(ctx, tenant)
		ctx = request.WithProject(ctx, project)
		render.RedirectWithRefresh(w, fmt.Sprintf("%s/%s/%d", routes.ProjectCtx(ctx), routes.ProjectSourceGithub, ghApp.UID))
	}
}

func HandleGithubInstall(tenantCtrl *tenant.Controller, projectCtrl *project.Controller, ghCtrl *githubapp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		installaltion_id := r.URL.Query().Get("installation_id")
		state := r.URL.Query().Get("source")
		setup_action := r.URL.Query().Get("setup_action")

		installID, err := strconv.ParseInt(installaltion_id, 10, 64)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error parsing installation id")
			render.Error500(w, r)
			return
		}

		ghAppState, err := parseGithubAppState(state)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error parsing github app uid from state")
			render.Error500(w, r)
			return
		}

		tenant, err := tenantCtrl.FindByUID(ctx, ghAppState.tenantUID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error finding tenant")
			render.Error500(w, r)
			return
		}

		project, err := projectCtrl.FindByUID(ctx, tenant.ID, ghAppState.projectUID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error finding tenant")
			render.Error500(w, r)
			return
		}

		ghApp, err := ghCtrl.FindByUID(ctx, tenant.ID, project.ID, ghAppState.ghAppUID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error getting github app")
			render.Error500(w, r)
			return
		}

		if ghApp == nil {
			log.Ctx(ctx).Error().Err(err).Msg("Github app not found")
			render.Error500(w, r)
			return
		}

		if setup_action != "setup" {
			err = ghCtrl.CompleteInstallation(ctx, ghApp, installID)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error completing github manifest")
				render.Error500(w, r)
				return
			}
		}

		ctx = request.WithTenant(ctx, tenant)
		ctx = request.WithProject(ctx, project)
		render.RedirectWithRefresh(w, fmt.Sprintf("%s/%s/%d", routes.ProjectCtx(ctx), routes.ProjectSourceGithub, ghApp.UID))
	}
}
