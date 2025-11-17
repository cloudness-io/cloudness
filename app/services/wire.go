package services

import (
	"github.com/cloudness-io/cloudness/app/services/cleanup"
	"github.com/cloudness-io/cloudness/job"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideServices,
)

type Services struct {
	JobScheduler *job.Scheduler
	Cleanup      *cleanup.Service
}

func ProvideServices(
	jobScheduler *job.Scheduler,
	cleanupSvc *cleanup.Service,
) Services {
	return Services{
		JobScheduler: jobScheduler,
		Cleanup:      cleanupSvc,
	}
}
