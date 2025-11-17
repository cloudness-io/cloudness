package pubsub

import "time"

type Provider string

const (
	ProviderMemory Provider = "inmemory"
	ProviderRedis  Provider = "redis"
)

type Config struct {
	App       string // app namespace prefix
	Namespace string

	Provider Provider

	HealthInterval time.Duration
	SendTimeout    time.Duration
	ChannelSize    int
}
