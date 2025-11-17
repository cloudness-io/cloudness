package gitpublic

import (
	"github.com/cloudness-io/cloudness/app/services/gitpublic"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(gitPublicSvc *gitpublic.Service) *Controller {
	return NewController(gitPublicSvc)
}
