package routes

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/types"
)

const (
	EnvironmentBase        = "environment"
	EnvironmentApplication = "/application"
	EnvironmentSettings    = "/settings"
	EnvironmentDelete      = "/delete"
	EnvironmentVolumes     = "/volumes"
	EnvironmentNav         = "/nav"
)

func Environment(uid int64) string {
	return fmt.Sprintf("%s/%d", EnvironmentBase, uid)
}

func EnvironmentNew() string {
	return fmt.Sprintf("%s/new", EnvironmentBase)
}

func EnvironmentCtx(ctx context.Context) string {
	environment, _ := request.EnvironmentFrom(ctx)
	return fmt.Sprintf("%s/environment/%d", ProjectCtx(ctx), environment.UID)
}

func EnvironmentCtxUID(ctx context.Context, envUID int64) string {
	return fmt.Sprintf("%s/environment/%d", ProjectCtx(ctx), envUID)
}

func EnvironmentUID(envUID int64) string {
	return fmt.Sprintf("%s/%d", EnvironmentBase, envUID)
}

func EnvironmentApplicationCtx(ctx context.Context) string {
	return fmt.Sprintf("%s/%s%s", EnvironmentCtx(ctx), EnvironmentApplication, TargetMainQuery)
}

func HasApplicationRoute(url string, env *types.Environment) bool {
	return strings.Contains(url, fmt.Sprintf("/environment/%d/application/", env.UID))
}
