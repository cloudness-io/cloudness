package canceler

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/cloudness-io/cloudness/app/pipeline/scheduler"
	"github.com/cloudness-io/cloudness/app/sse"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

type Canceler interface {
	CancelIncompleteBuilds(ctx context.Context, tenantID, projectID, applicationID int64) error
}

type canceler struct {
	deploymentStore store.DeploymentStore
	sseStreamer     sse.Streamer
	scheduler       scheduler.Scheduler
}

func New(
	deploymentStore store.DeploymentStore,
	sseStreamer sse.Streamer,
	scheduler scheduler.Scheduler,
) Canceler {
	return &canceler{
		deploymentStore: deploymentStore,
		sseStreamer:     sseStreamer,
		scheduler:       scheduler,
	}
}

func (c *canceler) CancelIncompleteBuilds(ctx context.Context, tenantID, projectID, applicationID int64) error {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	deployments, err := c.deploymentStore.ListIncompleteByApplicationID(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("canceler: failed to list incomplete deployments: %w", err)
	}

	for _, deployment := range deployments {
		now := time.Now().UTC().UnixMilli()
		deployment.Status = enum.DeploymentStatusCancelled
		deployment.Stopped = now
		if deployment.Started == 0 {
			deployment.Started = now
		}

		err := c.deploymentStore.Update(ctx, deployment)
		if err != nil {
			return fmt.Errorf("Canceler: failed to update deployment status to cancelled: %w", err)
		}

		if err := c.sseStreamer.Publish(ctx, projectID, enum.SSETypeDeploymentUpdated, deployment); err != nil {
			log.Debug().Err(err).Msg("canceler: failed to publish deployment update")
		}

		if err := c.scheduler.Cancel(ctx, deployment.ID); err != nil {
			log.Debug().Err(err).Msg("Canceler: failed to cancel scheduler")
		}

	}

	return nil
}
