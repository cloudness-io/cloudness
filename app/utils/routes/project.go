package routes

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/app/request"
)

const (
	ProjectBase        = "project"
	ProjectOverview    = "overview"
	ProjectEnvironment = "environment"
	ProjectConnections = "connections"
	ProjectMembers     = "members"
	ProjectSettings    = "settings"
	ProjectDelete      = "delete"
	ProjectNav         = "/nav"

	ProjectConnectionGithub = "connections/github"
)

func Project(uid int64) string {
	return fmt.Sprintf("%s/%d", ProjectBase, uid)
}

func ProjectNew() string {
	return fmt.Sprintf("%s/new", ProjectBase)
}

func ProjectNewCtx(ctx context.Context) string {
	return fmt.Sprintf("%s/%s", TenantCtx(ctx), ProjectNew())
}

func ProjectCtx(ctx context.Context) string {
	project, _ := request.ProjectFrom(ctx)
	return fmt.Sprintf("%s/project/%d", TenantCtx(ctx), project.UID)
}

func ProjectCtxUID(ctx context.Context, projectUID int64) string {
	return fmt.Sprintf("%s/project/%d", TenantCtx(ctx), projectUID)
}

func ProjectConnectionGithubUID(ghAppUID int64) string {
	return fmt.Sprintf("%s/%d", ProjectConnectionGithub, ghAppUID)
}

func ProjectConnectionGithubUIDCtx(ctx context.Context, ghAppUID int64) string {
	return fmt.Sprintf("%s/%s/%d", ProjectCtx(ctx), ProjectConnectionGithub, ghAppUID)
}
