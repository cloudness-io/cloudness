package job

import (
	"github.com/cloudness-io/cloudness/lock"
	"github.com/cloudness-io/cloudness/pubsub"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideExecutor,
	ProvideScheduler,
)

func ProvideExecutor(
	store Store,
	pubsubService pubsub.PubSub,
) *Executor {
	return NewExecutor(
		store,
		pubsubService,
	)
}

func ProvideScheduler(
	store Store,
	executor *Executor,
	mutexManager lock.MutexManager,
	pubsubService pubsub.PubSub,
) (*Scheduler, error) {
	return NewScheduler(
		store,
		executor,
		mutexManager,
		pubsubService,
	)
}
