package vcreate

import (
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/views/components/icons"
	"github.com/cloudness-io/cloudness/app/web/views/shared"
)

const (
	navGit      = "Git"
	navRegistry = "Registry"
	navDatabase = "Database"
	navOneclick = "OneClick"
)

func getCreatePageNavs() []*shared.PageNavItem {
	return []*shared.PageNavItem{
		{
			Name:      navGit,
			Icon:      icons.SourceGitIcon,
			ActionUrl: routes.AppNewGit,
		},
		{
			Name:      navRegistry,
			Icon:      icons.SourceRegistryIcon,
			ActionUrl: routes.AppNewRegistry,
		},
		{
			Name:      navDatabase,
			Icon:      icons.SourceDatabaseIcon,
			ActionUrl: routes.AppNewDatabase,
		},
		// {
		// 	Name:      navOneclick,
		// 	Icon:      icons.SourceOneclickIcon,
		// 	ActionUrl: routes.AppNewRegistry,
		// },
	}
}
