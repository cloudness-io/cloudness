package spec

import (
	"github.com/cloudness-io/cloudness/app/services/githubapp"
	"github.com/cloudness-io/cloudness/app/services/gitpublic"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideSpecService,
)

func ProvideSpecService(
	ghAppSvc *githubapp.Service,
	gitpublicSvc *gitpublic.Service,
) *Service {
	return NewService(
		ghAppSvc,
		gitpublicSvc,
	)
}
