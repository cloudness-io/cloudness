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
	serverStore  store.ServerStore
	tenantStore  store.TenantStore
	projectStore store.ProjectStore
	envStore     store.EnvironmentStore
	appStore     store.ApplicationStore
	volumeStore  store.VolumeStore
	tokenStore   store.TokenStore

	//factory
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
	tokenStore store.TokenStore,
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
		tokenStore:    tokenStore,
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

	if err := s.executor.Register(
		jobTypeTokens,
		newTokenCleanupJob(s.tokenStore),
	); err != nil {
		return fmt.Errorf("failed to register job handler for token cleanup: %w", err)
	}

	return nil
}

// scheduleRecurringCleanupJobs schedules the cleanup jobs.
func (s *Service) scheduleRecurringCleanupJobs(ctx context.Context) error {
	if err := s.scheduler.AddRecurring(
		ctx,
		jobTypeDeletedArtifacts,
		jobTypeDeletedArtifacts,
		jobCronDeletedArtifacts,
		jobMaxDurationDeletedArtifacts,
	); err != nil {
		return fmt.Errorf("failed to schedule deleted artifacts cleanup job: %w", err)
	}

	if err := s.scheduler.AddRecurring(
		ctx,
		jobTypeTokens,
		jobTypeTokens,
		jobCronTokens,
		jobMaxDurationTokens,
	); err != nil {
		return fmt.Errorf("failed to schedule token cleanup job: %w", err)
	}

	return nil
}
