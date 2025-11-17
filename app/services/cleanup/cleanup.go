package cleanup

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/app/services/manager"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/job"
)

type Service struct {
	//job related
	scheduler *job.Scheduler
	executor  *job.Executor

	//stores
	serverStore   store.ServerStore
	tenantStore   store.TenantStore
	projectStore  store.ProjectStore
	envStore      store.EnvironmentStore
	appStore      store.ApplicationStore
	volumeStore   store.VolumeStore
	serverFactory manager.ManagerFactory
}

func New(
	scheduler *job.Scheduler,
	executor *job.Executor,
	serverStore store.ServerStore,
	tenantStore store.TenantStore,
	projectStore store.ProjectStore,
	envStore store.EnvironmentStore,
	appStore store.ApplicationStore,
	volumeStore store.VolumeStore,
	serverFactory manager.ManagerFactory,
) *Service {
	return &Service{
		scheduler: scheduler,
		executor:  executor,

		serverStore:   serverStore,
		tenantStore:   tenantStore,
		projectStore:  projectStore,
		envStore:      envStore,
		appStore:      appStore,
		volumeStore:   volumeStore,
		serverFactory: serverFactory,
	}
}

func (s *Service) Register(ctx context.Context) error {
	if err := s.registerJobHandlers(); err != nil {
		return fmt.Errorf("failed to register cleanup job handlers: %w", err)
	}

	if err := s.scheduleRecurringCleanupJobs(ctx); err != nil {
		return fmt.Errorf("failed to schedule cleanup jobs: %w", err)
	}

	return nil
}

// registerJobHandlers registers handlers for all cleanup jobs.
func (s *Service) registerJobHandlers() error {
	if err := s.executor.Register(
		jobTypeDeletedArtifacts,
		newDeletedArtifactsJob(
			s.serverStore,
			s.tenantStore,
			s.projectStore,
			s.envStore,
			s.appStore,
			s.volumeStore,
			s.serverFactory,
		),
	); err != nil {
		return fmt.Errorf("failed to register job handler for deleted artifacts cleanup: %w", err)
	}
	return nil
}

// scheduleRecurringCleanupJobs schedules the cleanup jobs.
func (s *Service) scheduleRecurringCleanupJobs(ctx context.Context) error {
	err := s.scheduler.AddRecurring(
		ctx,
		jobTypeDeletedArtifacts,
		jobTypeDeletedArtifacts,
		jobCronDeletedArtifacts,
		jobMaxDurationDeletedArtifacts,
	)
	if err != nil {
		return fmt.Errorf("failed to schedule deleted artifacts cleanup job: %w", err)
	}
	return nil
}
