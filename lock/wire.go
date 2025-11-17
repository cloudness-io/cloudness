package lock

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideMutexManager,
)

func ProvideMutexManager(config Config) MutexManager {
	return NewInMemory(config)
}
