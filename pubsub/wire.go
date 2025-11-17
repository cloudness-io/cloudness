package pubsub

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

var WireSet = wire.NewSet(
	ProvidePubSub,
)

func ProvidePubSub(config Config, client redis.UniversalClient) PubSub {
	switch config.Provider {
	case ProviderRedis:
		return newRedis(client,
			WithApp(config.App),
			WithNamespace(config.Namespace),
			WithHealthCheckInterval(config.HealthInterval),
			WithSendTimeout(config.SendTimeout),
			WithSize(config.ChannelSize),
		)
	case ProviderMemory:
		fallthrough
	default:
		return newInMemory(
			WithApp(config.App),
			WithNamespace(config.Namespace),
			WithHealthCheckInterval(config.HealthInterval),
			WithSendTimeout(config.SendTimeout),
			WithSize(config.ChannelSize),
		)
	}
}
