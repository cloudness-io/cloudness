package githubapp

import (
	"github.com/cloudness-io/cloudness/app/services/githubapp"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(ghAppSvc *githubapp.Service) *Controller {
	return NewController(ghAppSvc)
}
