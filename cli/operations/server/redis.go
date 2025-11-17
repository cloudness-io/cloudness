package server

import (
	"github.com/cloudness-io/cloudness/types"

	"github.com/redis/go-redis/v9"
)

// ProvideRedis provides a redis client based on the configuration.
// TODO: add support for TLS.
func ProvideRedis(config *types.Config) (redis.UniversalClient, error) {
	options := &redis.Options{
		Addr:         config.Redis.Endpoint,
		MaxRetries:   config.Redis.MaxRetries,
		MinIdleConns: config.Redis.MinIdleConnections,
	}

	if config.Redis.Password != "" {
		options.Password = config.Redis.Password
	}
	if config.Redis.Username != "" {
		options.Username = config.Redis.Username
	}

	return redis.NewClient(options), nil
}
