package routes

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/types"
)

const (
	ProjectBase        = "project"
	ProjectOverview    = "overview"
	ProjectEnvironment = "environment"
	ProjectSource      = "source"
	ProjectMembers     = "members"
	ProjectSettings    = "settings"
	ProjectDelete      = "delete"
	ProjectNav         = "/nav"

	ProjectSourceGithub = "source/github"
)

func Project(uid int64) string {
	return fmt.Sprintf("%s/%d", ProjectBase, uid)
}

func ProjectCreate() string {
	return fmt.Sprintf("%s/create", ProjectBase)
}

func ProjectCtx(ctx context.Context) string {
	project, _ := request.ProjectFrom(ctx)
	return fmt.Sprintf("%s/project/%d", TenantCtx(ctx), project.UID)
}

func ProjectCtxUID(ctx context.Context, projectUID int64) string {
	return fmt.Sprintf("%s/project/%d", TenantCtx(ctx), projectUID)
}

func ProjectGithubCtx(ctx context.Context, ghApp *types.GithubApp) string {
	base := fmt.Sprintf("%s/%s", ProjectCtx(ctx), ProjectSourceGithub)
	if ghApp != nil {
		return fmt.Sprintf("%s/%d", base, ghApp.UID)
	}
	return base
}

func IsProjectUrl(url string) bool {
	return strings.Contains(url, "/project")
}
