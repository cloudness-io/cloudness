package volume

import (
	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/app/store"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(configSvc *config.Service, volumeStore store.VolumeStore) *Controller {
	return NewController(configSvc, volumeStore)
}
