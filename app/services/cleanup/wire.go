package cleanup

import (
	"github.com/cloudness-io/cloudness/app/services/manager"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/job"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideService,
)

func ProvideService(
	scheduler *job.Scheduler,
	executor *job.Executor,
	serverStore store.ServerStore,
	tenantStore store.TenantStore,
	projectStore store.ProjectStore,
	envStore store.EnvironmentStore,
	appStore store.ApplicationStore,
	volumeStore store.VolumeStore,
	tokenStore store.TokenStore,
	serverFactory manager.ManagerFactory,
) *Service {
	return New(
		scheduler,
		executor,
		serverStore,
		tenantStore,
		projectStore,
		envStore,
		appStore,
		volumeStore,
		tokenStore,
		serverFactory,
	)
}
