package dns

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideDNSService,
)

func ProvideDNSService() *Service {
	return New()
}
