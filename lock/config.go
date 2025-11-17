package lock

import "time"

type Provider string

const (
	MemoryProvider Provider = "inmemory"
	RedisProvider  Provider = "redis"
)

// A DelayFunc is used to decide the amount of time to wait between retries.
type DelayFunc func(tries int) time.Duration

type Config struct {
	App       string // app namespace prefix
	Namespace string
	Provider  Provider
	Expiry    time.Duration

	Tries      int
	RetryDelay time.Duration
	DelayFunc  DelayFunc

	DriftFactor   float64
	TimeoutFactor float64

	GenValueFunc func() (string, error)
	Value        string
}
