package manager

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

func (m *runnerManager) down(ctx context.Context, deployment *types.Deployment) error {
	log := log.With().
		Int64("deployment.id", deployment.ID).
		Int64("deployment.uid", deployment.UID).
		Int64("application.id", deployment.ApplicationID).
		Logger()

	app, err := m.applicationStore.Find(ctx, deployment.ApplicationID)
	if err != nil {
		log.Error().Err(err).Msg("manager: error finding application")
		return err
	}

	log = log.With().
		Str("application.name", app.Name).
		Logger()

	app.DeploymentStatus = deployment.Status.ToApplicationDeploymentStatus()
	if app.DeploymentStatus == enum.ApplicationDeploymentStatusSuccess {
		app.DeploymentID = &deployment.ID
	}
	if _, err := m.applicationStore.UpdateDeploymentStatus(ctx, app); err != nil {
		log.Error().Err(err).Msg("manager: error updating application")
		return err
	}

	if err := m.sseStreamer.Publish(ctx, app.ProjectID, enum.SSETypeDeploymentUpdated, deployment); err != nil {
		log.Warn().Err(err).Msg("manager: could not publish deployment updated event")
	}

	if err := m.sseStreamer.Publish(ctx, app.ProjectID, enum.SSETypeApplicationUpdated, app); err != nil {
		log.Warn().Err(err).Msg("manager: could not publish application updated event")
	}

	return nil
}
