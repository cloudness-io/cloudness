package manager

import (
	"github.com/cloudness-io/cloudness/app/controller/variable"
	"github.com/cloudness-io/cloudness/app/pipeline/scheduler"
	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/app/services/githubapp"
	"github.com/cloudness-io/cloudness/app/sse"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/logstream"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideRunnerManager,
)

func ProvideRunnerManager(
	scheduler scheduler.Scheduler,
	serverStore store.ServerStore,
	appStore store.ApplicationStore,
	deploymentStore store.DeploymentStore,
	volumeStore store.VolumeStore,
	logStore store.LogStore,
	varCtrl *variable.Controller,
	configSvc *config.Service,
	ghAppSvc *githubapp.Service,
	sseStreamer sse.Streamer,
	logStream logstream.LogStream,
) RunnerManager {
	return New(scheduler,
		serverStore,
		appStore,
		deploymentStore,
		volumeStore,
		logStore,
		varCtrl,
		configSvc,
		ghAppSvc,
		sseStreamer,
		logStream,
	)
}
