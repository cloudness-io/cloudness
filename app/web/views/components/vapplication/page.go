package vapplication

import (
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/views/components/icons"
	"github.com/cloudness-io/cloudness/app/web/views/shared"
	"github.com/cloudness-io/cloudness/types"
)

const (
	AppNavOverview    string = "Overview"
	AppNavDeployments string = "Deployments"
	AppNavMetrics     string = "Metrics"
	AppNavLogs        string = "Logs"
	AppNavTerminal    string = "Terminal"
	AppNavNetwork     string = "DNS"
	AppNavVolume      string = "Volumes"
	AppNavSettings    string = "Settings"
	AppNavEnvvars     string = "Variables"
	AppNavDelete      string = "Danger"
)

func getAppPageNav(app *types.Application) []*shared.PageNavItem {
	return []*shared.PageNavItem{
		{
			Name:      AppNavDeployments,
			Icon:      icons.DeploymentIcon,
			ActionUrl: routes.AppDeployments,
		},
		{
			Name:      AppNavMetrics,
			Icon:      icons.MetricsIcon,
			ActionUrl: routes.AppMetrics,
		},
		{
			Name:      AppNavLogs,
			Icon:      icons.LogsSectionIcon,
			ActionUrl: routes.AppLogs,
		},
		{
			Name:      AppNavTerminal,
			Icon:      icons.TerminalIcon,
			ActionUrl: routes.AppTerminal,
		},
		{
			Name:      AppNavVolume,
			Icon:      icons.VolumeSectionIcon,
			ActionUrl: routes.AppVolume,
		},
		{
			Name:      AppNavSettings,
			Icon:      icons.SettingsIcon,
			ActionUrl: routes.AppSettings,
		},
		{
			Name:      AppNavEnvvars,
			Icon:      icons.EnvironmentVariablesIcon,
			ActionUrl: routes.AppVariables,
		},
		{
			Name:      AppNavDelete,
			Icon:      icons.DeleteIcon,
			ActionUrl: routes.AppDelete,
		},
	}
}
