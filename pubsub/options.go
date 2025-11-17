package pubsub

import "time"

// An Option configures a pubsub instance.
type Option interface {
	Apply(*Config)
}

// OptionFunc is a function that configures a pubsub config.
type OptionFunc func(*Config)

// Apply calls f(config).
func (f OptionFunc) Apply(config *Config) {
	f(config)
}

// WithApp returns an option that set config app name.
func WithApp(value string) Option {
	return OptionFunc(func(m *Config) {
		m.App = value
	})
}

// WithNamespace returns an option that set config namespace.
func WithNamespace(value string) Option {
	return OptionFunc(func(m *Config) {
		m.Namespace = value
	})
}

// WithHealthCheckInterval specifies the config health check interval.
// PubSub will ping Server if it does not receive any messages
// within the interval (redis, ...).
// To disable health check, use zero interval.
func WithHealthCheckInterval(value time.Duration) Option {
	return OptionFunc(func(m *Config) {
		m.HealthInterval = value
	})
}

// WithSendTimeout specifies the pubsub send timeout after which
// the message is dropped.
func WithSendTimeout(value time.Duration) Option {
	return OptionFunc(func(m *Config) {
		m.SendTimeout = value
	})
}

// WithSize specifies the Go chan size in config that is used to buffer
// incoming messages.
func WithSize(value int) Option {
	return OptionFunc(func(m *Config) {
		m.ChannelSize = value
	})
}

type SubscribeConfig struct {
	topics         []string
	app            string
	namespace      string
	healthInterval time.Duration
	sendTimeout    time.Duration
	channelSize    int
}

type SubscribeOption interface {
	Apply(*SubscribeConfig)
}

// SubscribeOptionFunc is a function that configures a subscription config.
type SubscribeOptionFunc func(*SubscribeConfig)

// Apply calls f(subscribeConfig).
func (f SubscribeOptionFunc) Apply(config *SubscribeConfig) {
	f(config)
}

// WithNamespace returns an channel option that configures namespace.
func WithChannelNamespace(value string) SubscribeOption {
	return SubscribeOptionFunc(func(c *SubscribeConfig) {
		c.namespace = value
	})
}

type PublishConfig struct {
	app       string
	namespace string
}

type PublishOption interface {
	Apply(*PublishConfig)
}

// PublishOptionFunc is a function that configures a publish config.
type PublishOptionFunc func(*PublishConfig)

// Apply calls f(publishConfig).
func (f PublishOptionFunc) Apply(config *PublishConfig) {
	f(config)
}

// WithPublishNamespace modifies publish config namespace.
func WithPublishNamespace(value string) PublishOption {
	return PublishOptionFunc(func(c *PublishConfig) {
		c.namespace = value
	})
}

func formatTopic(app, ns, topic string) string {
	return app + ":" + ns + ":" + topic
}
