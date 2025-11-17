package scheduler

import (
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/lock"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideScheduler,
)

// ProvideScheduler provides a scheduler which can be used to schedule and request builds.
func ProvideScheduler(
	store store.DeploymentStore,
	lock lock.MutexManager,
) (Scheduler, error) {
	return newScheduler(store, lock)
}
