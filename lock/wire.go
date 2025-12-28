package lock

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

var WireSet = wire.NewSet(
	ProvideMutexManager,
)

func ProvideMutexManager(config Config, client redis.UniversalClient) MutexManager {
	switch config.Provider {
	case MemoryProvider:
		return NewInMemory(config)
	case RedisProvider:
		return NewRedis(config, client)
	}
	return NewInMemory(config)
}
