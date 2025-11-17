package vtenant

import (
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/views/components/icons"
	"github.com/cloudness-io/cloudness/app/web/views/shared"
	"github.com/cloudness-io/cloudness/types"
)

const (
	TenantNavSettings string = "Settings"
	TenantNavMembers  string = "Members"
	TenantNavLimits   string = "Limits"
	TenantNavDelete   string = "Danger"
)

func getTenantNav(tenant *types.Tenant) []*shared.PageNavItem {
	return []*shared.PageNavItem{
		{
			Name:      tenant.Name,
			Icon:      icons.TeamIcon,
			ActionUrl: routes.TenantSettings,
		},
		{
			Name:      TenantNavMembers,
			Icon:      icons.MembersIcon,
			ActionUrl: routes.TenantMembers,
		},
		{
			Name:      TenantNavLimits,
			Icon:      icons.LimitsIcon,
			ActionUrl: routes.TenantLimits,
		},
		{
			Name:      TenantNavDelete,
			Icon:      icons.DeleteIcon,
			ActionUrl: routes.TenantDelete,
		},
	}
}
