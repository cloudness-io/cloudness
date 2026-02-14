package vproject

import (
	"context"

	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/views/components/icons"
	"github.com/cloudness-io/cloudness/app/web/views/shared"
)

const (
	ProjectNavOverview     string = "Overview"
	ProjectNavEnvironments string = "Environments"
	ProjectNavConnections  string = "Connections"
	ProjectNavMembers      string = "Members"
	ProjectNavSettings     string = "Settings"
)

func getProjectPageNav(ctx context.Context) []*shared.PageNavItem {
	return []*shared.PageNavItem{
		{
			Name:      ProjectNavOverview,
			Icon:      icons.ProjectIcon,
			ActionUrl: routes.ProjectOverview,
		},
		{
			Name:      ProjectNavEnvironments,
			Icon:      icons.EnvironmentIcon,
			ActionUrl: routes.ProjectEnvironment,
		},
		{
			Name:      ProjectNavConnections,
			Icon:      icons.SourceGitIcon,
			ActionUrl: routes.ProjectConnections,
		},

		{
			Name:      ProjectNavMembers,
			Icon:      icons.MembersIcon,
			ActionUrl: routes.ProjectMembers,
			Hide:      !request.IsProjectOwner(ctx),
		},
		{
			Name:      ProjectNavSettings,
			Icon:      icons.SettingsIcon,
			ActionUrl: routes.ProjectSettings,
		},
	}
}
