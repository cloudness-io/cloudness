package tenant

import (
	"context"

	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/types"
)

func canEdit(ctx context.Context, tenantCtrl *tenant.Controller, tenant *types.Tenant) bool {
	restrictions := tenantCtrl.GetRestrctions(ctx, tenant)

	if request.IsSuperAdmin(ctx) || (restrictions.AllowAdminToModify && request.IsTeamAdmin(ctx)) {
		return true
	}
	return false
}
