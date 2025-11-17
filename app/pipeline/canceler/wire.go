package canceler

import (
	"github.com/cloudness-io/cloudness/app/pipeline/scheduler"
	"github.com/cloudness-io/cloudness/app/sse"
	"github.com/cloudness-io/cloudness/app/store"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideCanceler,
)

func ProvideCanceler(
	deploymentStore store.DeploymentStore,
	sseStreamer sse.Streamer,
	scheduler scheduler.Scheduler,
) Canceler {
	return New(deploymentStore, sseStreamer, scheduler)
}
