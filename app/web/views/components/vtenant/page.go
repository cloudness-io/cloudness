package vtenant

import (
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/views/components/icons"
	"github.com/cloudness-io/cloudness/app/web/views/shared"
	"github.com/cloudness-io/cloudness/types"
)

const (
	TenantNavProjects     string = "Projects"
	TenantNavSettings     string = "Team Settings"
	TenantNavMembers      string = "Team"
	TenantNavRestrictions string = "Restrictions"
	TenantNavDelete       string = "Danger"
)

func getTenantNav(tenant *types.Tenant, canEdit bool) []*shared.PageNavItem {
	return []*shared.PageNavItem{
		{
			Name:      TenantNavProjects,
			Icon:      icons.HomeIcon,
			ActionUrl: routes.TenantUID(tenant.UID),
		},
		{
			Name:      TenantNavMembers,
			Icon:      icons.TeamMembersIcon,
			ActionUrl: routes.TenantMembers,
			Disabled:  !canEdit,
		},
		{
			Name:      TenantNavRestrictions,
			Icon:      icons.LimitsIcon,
			ActionUrl: routes.TenantRestrictions,
			Disabled:  !canEdit,
			Hide:      !canEdit,
		},
		{
			Name:      TenantNavSettings,
			Icon:      icons.SettingsIcon,
			ActionUrl: routes.TenantSettings,
			Disabled:  !canEdit,
		},
	}
}
