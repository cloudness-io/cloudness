package vapplication

import (
	"fmt"

	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/views/components/icons"
	"github.com/cloudness-io/cloudness/app/web/views/shared"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
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
			ActionUrl: fmt.Sprintf("%s/%s", routes.AppMetrics, enum.MetricsSpan1h),
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
			Name:      AppNavEnvvars,
			Icon:      icons.EnvironmentVariablesIcon,
			ActionUrl: routes.AppVariables,
		},
		{
			Name:      AppNavSettings,
			Icon:      icons.SettingsIcon,
			ActionUrl: routes.AppSettings,
		},
		{
			Name:      AppNavDelete,
			Icon:      icons.DeleteIcon,
			ActionUrl: routes.AppDelete,
		},
	}
}
