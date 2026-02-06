package routes

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/app/request"
)

const (
	TenantBase = "team"

	TenantSettings      = "settings"
	TenantMembers       = "members"
	TenantRestrictions  = "restrictions"
	TenantDelete        = "delete"
	TenantMembersAction = "/members"
)

func TenantBaseURL() string {
	return "/" + TenantBase
}

func TenantUID(uid int64) string {
	return fmt.Sprintf("/%s/%d", TenantBase, uid)
}

func TenantNew() string {
	return fmt.Sprintf("/%s/new", TenantBase)
}

func TenantCtx(ctx context.Context) string {
	tenant, ok := request.TenantFrom(ctx)
	if !ok {
		return "/"
	}
	return fmt.Sprintf("/team/%d", tenant.UID)
}

func TenantMembersUrl(ctx context.Context) string {
	return fmt.Sprintf("%s%s", TenantCtx(ctx), TenantMembersAction)
}
