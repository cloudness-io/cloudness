package proxy

import "github.com/google/wire"

var WireSet = wire.NewSet(
	ProvideProxyService,
)

func ProvideProxyService() *Service {
	return New()
}
