package vaccount

import (
	"github.com/cloudness-io/cloudness/app/utils/routes"
	"github.com/cloudness-io/cloudness/app/web/views/components/icons"
	"github.com/cloudness-io/cloudness/app/web/views/shared"
)

const (
	AccountNavInfo    string = "Account"
	AccountNavSession string = "Sessions"
	AccountNavTeams   string = "Teams"
	AccountNavDelete  string = "Delete"
)

func getAccountPageNav() []*shared.PageNavItem {
	return []*shared.PageNavItem{
		{
			Name:      AccountNavInfo,
			Icon:      icons.UserProfileIcon,
			ActionUrl: routes.AccountProfile,
		},
		{
			Name:      AccountNavSession,
			Icon:      icons.UserSessionIcon,
			ActionUrl: routes.AccountSession,
		},
		{
			Name:      AccountNavTeams,
			Icon:      icons.TeamIcon,
			ActionUrl: routes.AccountTeams,
		},
		{
			Name:      AccountNavDelete,
			Icon:      icons.DeleteIcon,
			ActionUrl: routes.AccountDelete,
		},
	}
}
