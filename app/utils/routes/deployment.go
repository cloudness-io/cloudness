package routes

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/app/request"
)

const (
	DeploymentBase = "deployment"
)

func Deployment(uid int64) string {
	return fmt.Sprintf("%s/%d", DeploymentBase, uid)
}

func DeploymentLog(uid int64) string {
	return fmt.Sprintf("%s/logs", Deployment(uid))
}

func DeploymentLogStream(uid int64) string {
	return fmt.Sprintf("%s/logs/stream", Deployment(uid))
}

func DeploymentCtx(ctx context.Context) string {
	deployment, _ := request.DeploymentFrom(ctx)
	return fmt.Sprintf("%s/deployment/%d", ApplicationCtx(ctx), deployment.UID)
}
