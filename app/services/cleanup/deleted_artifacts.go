package cleanup

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudness-io/cloudness/app/services/manager"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/job"

	"github.com/rs/zerolog/log"
)

const (
	jobTypeDeletedArtifacts        = "cloudness:cleanup:deleted-artifacts"
	jobCronDeletedArtifacts        = "*/5 * * * *"
	jobMaxDurationDeletedArtifacts = 10 * time.Minute
	maxRetrivalItems               = 100
)

type deletedArtifactsJob struct {
	serverStore  store.ServerStore
	tenantStore  store.TenantStore
	projectStore store.ProjectStore
	envStore     store.EnvironmentStore
	appStore     store.ApplicationStore
	volumeStore  store.VolumeStore

	serverFactory manager.ManagerFactory
}

func newDeletedArtifactsJob(
	serverStore store.ServerStore,
	tenantStore store.TenantStore,
	projectStore store.ProjectStore,
	envStore store.EnvironmentStore,
	appStore store.ApplicationStore,
	volumeStore store.VolumeStore,
	serverFactory manager.ManagerFactory,
) *deletedArtifactsJob {
	return &deletedArtifactsJob{
		serverStore:   serverStore,
		tenantStore:   tenantStore,
		projectStore:  projectStore,
		envStore:      envStore,
		appStore:      appStore,
		volumeStore:   volumeStore,
		serverFactory: serverFactory,
	}
}

type purgeTracker struct {
	tenant  int64
	project int64
	env     int64
	app     int64
	volume  int64
}

func (t *purgeTracker) ToString() string {
	return fmt.Sprintf("Artifacts Purge successful: %d tenants, %d projects, %d environments, %d apps, %d volumes", t.tenant, t.project, t.env, t.app, t.volume)
}

func (j *deletedArtifactsJob) Handle(ctx context.Context, _ string, _ job.ProgressReporter) (string, error) {
	olderThan := time.Now().UTC().Add(-time.Minute)

	log.Ctx(ctx).Info().Msgf("cleanup: starting to purge artifacts deleted before %s: %d", olderThan.Format(time.RFC3339Nano), olderThan.UnixMilli())

	deletedBeforeOrAt := olderThan.UnixMilli()

	tracker := &purgeTracker{}

	if err := j.cleanupApps(ctx, nil, nil, nil, &deletedBeforeOrAt, tracker); err != nil {
		return "Error cleaning up application", err
	}

	if err := j.cleanupVolumes(ctx, nil, nil, nil, &deletedBeforeOrAt, tracker); err != nil {
		return "Error cleaning up volume", err
	}

	if err := j.cleanupEnviornments(ctx, nil, nil, &deletedBeforeOrAt, tracker); err != nil {
		return "Error cleaning up environment", err
	}

	if err := j.cleanupProjects(ctx, nil, &deletedBeforeOrAt, tracker); err != nil {
		return "Error cleaning up project", err
	}

	if err := j.cleanupTenant(ctx, &deletedBeforeOrAt, tracker); err != nil {
		return "Error cleaning up tenant", err
	}

	return tracker.ToString(), nil
}
