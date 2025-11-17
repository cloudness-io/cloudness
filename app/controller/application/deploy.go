package application

import (
	"context"

	"github.com/cloudness-io/cloudness/app/pipeline/triggerer"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

func (c *Controller) NeedsDeployment(ctx context.Context, application *types.Application) (*types.Application, error) {
	return c.applicationStore.UpdateNeedsDeployment(ctx, application)
}

func (c *Controller) Redeploy(ctx context.Context, actor string, application *types.Application) (*types.Deployment, error) {
	title, err := c.specSvc.GetCurrentTag(ctx, application)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Error finding deploy trigger message")
		title = "Manual Trigger"
	}

	hook := &triggerer.TriggerHook{
		ApplicaitonID: application.ID,
		Triggerer:     actor,
		Title:         title,
		Action:        enum.TriggerActionManual,
	}

	return c.triggerer.Trigger(ctx, hook)
}

func (c *Controller) Deploy(ctx context.Context, actor string, action enum.TriggerAction, application *types.Application) (*types.Deployment, error) {
	triggerMsg := "Deployment triggered"
	tag, err := c.specSvc.GetLatestTag(ctx, application)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Error updating application tags like git commit and trigger message")
		// Default tag messages
		switch action {
		case enum.TriggerActionCreate:
			triggerMsg = "Application created"
		case enum.TriggerActionHook:
			triggerMsg = "Application hook trigger"
		case enum.TriggerActionManual:
			triggerMsg = "Application updated"
		}
	}

	triggerMsg = tag

	hook := &triggerer.TriggerHook{
		ApplicaitonID: application.ID,
		Triggerer:     actor,
		Title:         triggerMsg,
		Action:        action,
	}

	return c.triggerer.Trigger(ctx, hook)
}
