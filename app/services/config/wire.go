package config

import (
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/types"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideConfigService,
)

func ProvideConfigService(config *types.Config, instanceStore store.InstanceStore) *Service {
	return NewService(config, instanceStore)
}
