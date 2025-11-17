package gitpublic

import "github.com/google/wire"

var WireSet = wire.NewSet(ProvideGitpublicService)

func ProvideGitpublicService() *Service {
	return NewService()
}
