package routes

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/app/request"
)

const (
	ApplicationBase = "application"

	AppNewGit       = "application/new/git"
	AppNewGithub    = "application/new/git/github"
	AppNewGitPublic = "application/new/git/git-public"
	AppNewRegistry  = "application/new/registry"
	AppNewDatabase  = "application/new/database"
	AppNewOneclick  = "application/new/oneclick"

	AppDeployments         = "deployments"
	AppMetrics             = "metrics"
	AppLogs                = "logs"
	AppTerminal            = "terminal"
	AppSource              = "source"
	AppVolume              = "volumes"
	AppVolumeDetach        = "detach"
	AppVolumeUnAttached    = "volumes/unattached"
	AppVolumeCreate        = "volumes/create"
	AppSettings            = "settings"
	AppVariables           = "variables"
	AppDelete              = "delete"
	AppNetworkPrivate      = "network/private"
	AppNetworkHTTP         = "network/http"
	AppNetworkHTTPGenerate = AppNetworkHTTP + "/generate"
	AppNetworkTCP          = "network/tcp"
	AppFavorite            = "favorite"

	AppNav = "/nav"

	//General Header
	AppIcon = "icon"
	AppName = "name"

	AppSourceGithub    = "source/github"
	AppSourceGitPublic = "source/git-public"
	AppSourceRegistry  = "source/registry"
)

func Application(envUID, appUID int64) string {
	return fmt.Sprintf("environment/%d/%s/%d", envUID, ApplicationBase, appUID)
}

func ApplicationNew(envUID int64) string {
	return fmt.Sprintf("environment/%d/%s/new/git", envUID, ApplicationBase)
}

func ApplicationCtx(ctx context.Context) string {
	app, _ := request.ApplicationFrom(ctx)
	return fmt.Sprintf("%s/application/%d", EnvironmentCtx(ctx), app.UID)
}

func ApplicationCtxUID(ctx context.Context, appUID int64) string {
	return fmt.Sprintf("%s/application/%d", EnvironmentCtx(ctx), appUID)
}

func ApplicationProject(projectUID, envUID, appUID int64) string {
	return fmt.Sprintf("%s/%s", Project(projectUID), Application(envUID, appUID))
}
