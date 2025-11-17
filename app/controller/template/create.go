package template

import (
	"context"

	"github.com/cloudness-io/cloudness/app/auth"
	"github.com/cloudness-io/cloudness/app/pipeline/triggerer"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

func (c *Controller) Create(
	ctx context.Context,
	session *auth.Session,
	tenant *types.Tenant,
	project *types.Project,
	env *types.Environment,
	template *types.Template,
) error {
	apps, err := c.appCtrl.CreateFromTemplate(ctx, session.Principal.DisplayName, tenant, project, env, template)
	if err != nil {
		return err
	}

	for _, app := range apps {
		hook := &triggerer.TriggerHook{
			ApplicaitonID: app.ID,
			Triggerer:     session.Principal.DisplayName,
			Title:         "Created from template",
			Action:        enum.TriggerActionManual,
		}

		if _, err := c.triggerer.Trigger(ctx, hook); err != nil {
			//Intentionally avoid error handling here
			log.Ctx(ctx).Error().Err(err).Msg("Error triggering deployment during template creations")
		}
	}
	return nil
}
